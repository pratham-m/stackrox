package alertmanager

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	alertMocks "github.com/stackrox/rox/central/alert/datastore/mocks"
	notifierMocks "github.com/stackrox/rox/central/notifier/processor/mocks"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/violationmessages/printer"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	nowProcess        = getProcessIndicator(ptypes.TimestampNow())
	yesterdayProcess  = getProcessIndicator(protoconv.ConvertTimeToTimestamp(time.Now().Add(-24 * time.Hour)))
	twoDaysAgoProcess = getProcessIndicator(protoconv.ConvertTimeToTimestamp(time.Now().Add(-2 * 24 * time.Hour)))

	firstKubeEventViolation  = getKubeEventViolation("1", protoconv.ConvertTimeToTimestamp(time.Now().Add(-24*time.Hour)))
	secondKubeEventViolation = getKubeEventViolation("2", ptypes.TimestampNow())

	firstNetworkFlowViolation  = getNetworkFlowViolation("1", protoconv.ConvertTimeToTimestamp(time.Now().Add(-24*time.Hour)))
	secondNetworkFlowViolation = getNetworkFlowViolation("2", ptypes.TimestampNow())
)

func getKubeEventViolation(msg string, timestamp *ptypes.Timestamp) *storage.Alert_Violation {
	return &storage.Alert_Violation{
		Message: msg,
		Type:    storage.Alert_Violation_K8S_EVENT,
		Time:    timestamp,
	}
}

func getNetworkFlowViolation(msg string, networkFlowTimestamp *ptypes.Timestamp) *storage.Alert_Violation {
	return &storage.Alert_Violation{
		Message: msg,
		MessageAttributes: &storage.Alert_Violation_KeyValueAttrs_{
			KeyValueAttrs: &storage.Alert_Violation_KeyValueAttrs{
				Attrs: []*storage.Alert_Violation_KeyValueAttrs_KeyValueAttr{
					{
						Key: "NetworkFlowTimestamp",
						Value: protoconv.
							ConvertTimestampToTimeOrNow(networkFlowTimestamp).
							Format("2006-01-02 15:04:05 UTC"),
					},
				},
			},
		},
		Type: storage.Alert_Violation_NETWORK_FLOW,
	}
}

func getProcessIndicator(timestamp *ptypes.Timestamp) *storage.ProcessIndicator {
	return &storage.ProcessIndicator{
		Signal: &storage.ProcessSignal{
			Name: "apt-get",
			Time: timestamp,
		},
	}
}

func getFakeRuntimeAlert(indicators ...*storage.ProcessIndicator) *storage.Alert {
	v := &storage.Alert_ProcessViolation{Processes: indicators}
	printer.UpdateProcessAlertViolationMessage(v)
	return &storage.Alert{
		LifecycleStage:   storage.LifecycleStage_RUNTIME,
		ProcessViolation: v,
	}
}

func appendViolations(alert *storage.Alert, violations ...*storage.Alert_Violation) *storage.Alert {
	alert.Violations = append(alert.Violations, violations...)
	return alert
}

func TestAlertManager(t *testing.T) {
	suite.Run(t, new(AlertManagerTestSuite))
}

type AlertManagerTestSuite struct {
	suite.Suite

	alertsMock   *alertMocks.MockDataStore
	notifierMock *notifierMocks.MockProcessor

	alertManager AlertManager

	mockCtrl *gomock.Controller
	ctx      context.Context

	envIsolator *envisolator.EnvIsolator
}

func (suite *AlertManagerTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.alertsMock = alertMocks.NewMockDataStore(suite.mockCtrl)
	suite.notifierMock = notifierMocks.NewMockProcessor(suite.mockCtrl)

	suite.alertManager = New(suite.notifierMock, suite.alertsMock, nil)
	suite.envIsolator = envisolator.NewEnvIsolator(suite.T())
}

func (suite *AlertManagerTestSuite) TearDownTest() {
	suite.envIsolator.RestoreAll()
	suite.mockCtrl.Finish()
}

// Returns a function that can be used to match *v1.Query,
// which ensure that the query specifies all the fields.
func queryHasFields(fields ...search.FieldLabel) func(interface{}) bool {
	return func(in interface{}) bool {
		q := in.(*v1.Query)

		fieldsFound := make([]bool, len(fields))
		search.ApplyFnToAllBaseQueries(q, func(bq *v1.BaseQuery) {
			mfQ, ok := bq.GetQuery().(*v1.BaseQuery_MatchFieldQuery)
			if !ok {
				return
			}
			for i, field := range fields {
				if mfQ.MatchFieldQuery.GetField() == field.String() {
					fieldsFound[i] = true
				}
			}
		})

		for _, found := range fieldsFound {
			if !found {
				return false
			}
		}
		return true
	}
}

func (suite *AlertManagerTestSuite) TestNotifyAndUpdateBatch() {
	alerts := []*storage.Alert{fixtures.GetAlert(), fixtures.GetAlert()}
	alerts[0].GetPolicy().Id = "Pol1"
	alerts[0].GetDeployment().Id = "Dep1"
	alerts[1].GetPolicy().Id = "Pol2"
	alerts[1].GetDeployment().Id = "Dep2"

	envIsolator := envisolator.NewEnvIsolator(suite.T())
	defer envIsolator.RestoreAll()
	envIsolator.Setenv(env.AlertRenotifDebounceDuration.EnvVar(), "5m")

	resolvedAlerts := []*storage.Alert{alerts[0].Clone(), alerts[1].Clone()}
	resolvedAlerts[0].ResolvedAt = protoconv.MustConvertTimeToTimestamp(time.Now().Add(-10 * time.Minute))
	resolvedAlerts[1].ResolvedAt = protoconv.MustConvertTimeToTimestamp(time.Now().Add(-2 * time.Minute))

	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx,
		testutils.PredMatcher("query for dep 1", func(q *v1.Query) bool {
			return strings.Contains(proto.MarshalTextString(q), "Dep1")
		})).Return([]*storage.Alert{resolvedAlerts[0]}, nil)
	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx,
		testutils.PredMatcher("query for dep 2", func(q *v1.Query) bool {
			return strings.Contains(proto.MarshalTextString(q), "Dep2")
		})).Return([]*storage.Alert{resolvedAlerts[1]}, nil)

	// Only the first alert will get notified
	suite.notifierMock.EXPECT().ProcessAlert(suite.ctx, alerts[0])
	// All alerts will still get inserted
	for _, alert := range alerts {
		suite.alertsMock.EXPECT().UpsertAlert(suite.ctx, alert)
	}
	suite.NoError(suite.alertManager.(*alertManagerImpl).notifyAndUpdateBatch(suite.ctx, alerts))
}

func (suite *AlertManagerTestSuite) TestGetAlertsByPolicy() {
	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx, testutils.PredMatcher("query for violation state, policy", queryHasFields(search.ViolationState, search.PolicyID))).Return(([]*storage.Alert)(nil), nil)

	modified, err := suite.alertManager.AlertAndNotify(suite.ctx, nil, WithPolicyID("pid"))
	suite.False(modified.Cardinality() > 0)
	suite.NoError(err, "update should succeed")
}

func (suite *AlertManagerTestSuite) TestGetAlertsByDeployment() {
	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx, testutils.PredMatcher("query for violation state, deployment", queryHasFields(search.ViolationState, search.DeploymentID))).Return(([]*storage.Alert)(nil), nil)

	modified, err := suite.alertManager.AlertAndNotify(suite.ctx, nil, WithDeploymentID("did", false))
	suite.False(modified.Cardinality() > 0)
	suite.NoError(err, "update should succeed")
}

func (suite *AlertManagerTestSuite) TestOnUpdatesWhenAlertsDoNotChange() {
	alerts := getAlerts()

	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx, gomock.Any()).Return(alerts, nil)
	// No updates should be attempted

	modified, err := suite.alertManager.AlertAndNotify(suite.ctx, alerts)
	suite.False(modified.Cardinality() > 0)
	suite.NoError(err, "update should succeed")
}

func (suite *AlertManagerTestSuite) TestMarksOldAlertsStale() {
	alerts := getAlerts()

	suite.alertsMock.EXPECT().MarkAlertStale(suite.ctx, alerts[0].GetId()).Return(nil)

	// Unchanged alerts should not be updated.

	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx, gomock.Any()).Return(alerts, nil)
	// We should get a notification for the new alert.
	suite.notifierMock.EXPECT().ProcessAlert(gomock.Any(), alerts[0]).Return()

	// Make one of the alerts not appear in the current alerts.
	modified, err := suite.alertManager.AlertAndNotify(suite.ctx, alerts[1:])
	suite.True(modified.Cardinality() > 0)
	suite.NoError(err, "update should succeed")
}

func (suite *AlertManagerTestSuite) TestSendsNotificationsForNewAlerts() {
	alerts := getAlerts()

	// Only the new alert will be updated.
	suite.alertsMock.EXPECT().UpsertAlert(suite.ctx, alerts[0]).Return(nil)

	// We should get a notification for the new alert.
	suite.notifierMock.EXPECT().ProcessAlert(gomock.Any(), alerts[0]).Return()

	// Make one of the alerts not appear in the previous alerts.
	suite.alertsMock.EXPECT().SearchRawAlerts(suite.ctx, gomock.Any()).Return(alerts[1:], nil)

	modified, err := suite.alertManager.AlertAndNotify(suite.ctx, alerts)
	suite.True(modified.Cardinality() > 0)
	suite.NoError(err, "update should succeed")
}

func TestMergeProcessesFromOldIntoNew(t *testing.T) {
	for _, c := range []struct {
		desc           string
		old            *storage.Alert
		new            *storage.Alert
		expectedNew    *storage.Alert
		expectedOutput bool
	}{
		{
			desc:           "Equal",
			old:            getFakeRuntimeAlert(yesterdayProcess),
			new:            getFakeRuntimeAlert(yesterdayProcess),
			expectedNew:    nil,
			expectedOutput: false,
		},
		{
			desc:           "Equal with two",
			old:            getFakeRuntimeAlert(yesterdayProcess, nowProcess),
			new:            getFakeRuntimeAlert(yesterdayProcess, nowProcess),
			expectedOutput: false,
		},
		{
			desc:           "New has new",
			old:            getFakeRuntimeAlert(yesterdayProcess),
			new:            getFakeRuntimeAlert(nowProcess),
			expectedNew:    getFakeRuntimeAlert(yesterdayProcess, nowProcess),
			expectedOutput: true,
		},
		{
			desc:           "New has many new",
			old:            getFakeRuntimeAlert(twoDaysAgoProcess, yesterdayProcess),
			new:            getFakeRuntimeAlert(yesterdayProcess, nowProcess),
			expectedNew:    getFakeRuntimeAlert(twoDaysAgoProcess, yesterdayProcess, nowProcess),
			expectedOutput: true,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			out := mergeProcessesFromOldIntoNew(c.old, c.new)
			assert.Equal(t, c.expectedOutput, out)
			if c.expectedNew != nil {
				assert.Equal(t, c.expectedNew, c.new)
			}
		})
	}
}

func TestMergeRunTimeAlerts(t *testing.T) {
	for _, c := range []struct {
		desc           string
		old            *storage.Alert
		new            *storage.Alert
		expectedNew    *storage.Alert
		expectedOutput bool
	}{
		{
			desc:           "No process; no event",
			old:            getFakeRuntimeAlert(),
			new:            getFakeRuntimeAlert(),
			expectedOutput: false,
		},
		{
			desc:           "No new process; no event",
			old:            getFakeRuntimeAlert(yesterdayProcess),
			new:            getFakeRuntimeAlert(),
			expectedOutput: false,
		},
		{
			desc:           "No process; no new event",
			old:            appendViolations(getFakeRuntimeAlert(), firstKubeEventViolation),
			new:            getFakeRuntimeAlert(),
			expectedOutput: false,
		},
		{
			desc:           "No process; new event",
			old:            getFakeRuntimeAlert(),
			new:            appendViolations(getFakeRuntimeAlert(), firstKubeEventViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(), firstKubeEventViolation),
			expectedOutput: true,
		},
		{
			desc:           "Equal process; no new event",
			old:            appendViolations(getFakeRuntimeAlert(yesterdayProcess), firstKubeEventViolation),
			new:            appendViolations(getFakeRuntimeAlert(yesterdayProcess)),
			expectedOutput: false,
		},
		{
			desc:           "Equal process; new event",
			old:            appendViolations(getFakeRuntimeAlert(yesterdayProcess), firstKubeEventViolation),
			new:            appendViolations(getFakeRuntimeAlert(yesterdayProcess), secondKubeEventViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(yesterdayProcess), secondKubeEventViolation, firstKubeEventViolation),
			expectedOutput: true,
		},
		{
			desc:           "New process; new event ",
			old:            appendViolations(getFakeRuntimeAlert(yesterdayProcess), firstKubeEventViolation),
			new:            appendViolations(getFakeRuntimeAlert(nowProcess), secondKubeEventViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(yesterdayProcess, nowProcess), secondKubeEventViolation, firstKubeEventViolation),
			expectedOutput: true,
		},
		{
			desc:           "New process; no new event ",
			old:            appendViolations(getFakeRuntimeAlert(yesterdayProcess), firstKubeEventViolation),
			new:            getFakeRuntimeAlert(nowProcess),
			expectedNew:    getFakeRuntimeAlert(yesterdayProcess, nowProcess),
			expectedOutput: true,
		},
		{
			desc:           "Many new process; many new events",
			old:            getFakeRuntimeAlert(twoDaysAgoProcess, yesterdayProcess),
			new:            appendViolations(getFakeRuntimeAlert(yesterdayProcess, nowProcess), firstKubeEventViolation, secondKubeEventViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(twoDaysAgoProcess, yesterdayProcess, nowProcess), firstKubeEventViolation, secondKubeEventViolation),
			expectedOutput: true,
		},
		{
			desc:           "No process; new network flow",
			old:            getFakeRuntimeAlert(),
			new:            appendViolations(getFakeRuntimeAlert(), firstNetworkFlowViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(), firstNetworkFlowViolation),
			expectedOutput: true,
		},
		{
			desc:           "Old process with old flow; new network flow",
			old:            appendViolations(getFakeRuntimeAlert(nowProcess), firstNetworkFlowViolation),
			new:            appendViolations(getFakeRuntimeAlert(nowProcess), secondNetworkFlowViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(nowProcess), secondNetworkFlowViolation, firstNetworkFlowViolation),
			expectedOutput: true,
		},
		{
			desc:           "Many new process; many new flow",
			old:            appendViolations(getFakeRuntimeAlert(twoDaysAgoProcess)),
			new:            appendViolations(getFakeRuntimeAlert(yesterdayProcess, nowProcess), firstNetworkFlowViolation, secondNetworkFlowViolation),
			expectedNew:    appendViolations(getFakeRuntimeAlert(twoDaysAgoProcess, yesterdayProcess, nowProcess), firstNetworkFlowViolation, secondNetworkFlowViolation),
			expectedOutput: true,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			out := mergeRunTimeAlerts(c.old, c.new)
			assert.Equal(t, c.expectedOutput, out)
			if c.expectedNew != nil {
				assert.Equal(t, c.expectedNew, c.new)
			}
		})
	}
}

//////////////////////////////////////
// TEST DATA
///////////////////////////////////////

// Policies are set up so that policy one is violated by deployment 1, 2 is violated by 2, etc.
func getAlerts() []*storage.Alert {
	return []*storage.Alert{
		{
			Id:     "alert1",
			Policy: getPolicies()[0],
			Entity: &storage.Alert_Deployment_{Deployment: getDeployments()[0]},
			Time:   &ptypes.Timestamp{Seconds: 100},
		},
		{
			Id:     "alert2",
			Policy: getPolicies()[1],
			Entity: &storage.Alert_Deployment_{Deployment: getDeployments()[1]},
			Time:   &ptypes.Timestamp{Seconds: 200},
		},
		{
			Id:     "alert3",
			Policy: getPolicies()[2],
			Entity: &storage.Alert_Deployment_{Deployment: getDeployments()[2]},
			Time:   &ptypes.Timestamp{Seconds: 300},
		},
	}
}

// Policies are set up so that policy one is violated by deployment 1, 2 is violated by 2, etc.
func getDeployments() []*storage.Alert_Deployment {
	return []*storage.Alert_Deployment{
		{
			Name: "deployment1",
			Containers: []*storage.Alert_Deployment_Container{
				{
					Image: &storage.ContainerImage{
						Name: &storage.ImageName{
							Tag:    "latest1",
							Remote: "stackrox/health",
						},
					},
				},
			},
		},
		{
			Name: "deployment2",
			Containers: []*storage.Alert_Deployment_Container{
				{
					Image: &storage.ContainerImage{
						Name: &storage.ImageName{
							Tag:    "latest2",
							Remote: "stackrox/health",
						},
					},
				},
			},
		},
		{
			Name: "deployment3",
			Containers: []*storage.Alert_Deployment_Container{
				{
					Image: &storage.ContainerImage{
						Name: &storage.ImageName{
							Tag:    "latest3",
							Remote: "stackrox/health",
						},
					},
				},
			},
		},
	}
}

// Policies are set up so that policy one is violated by deployment 1, 2 is violated by 2, etc.
func getPolicies() []*storage.Policy {
	return []*storage.Policy{
		{
			Id:         "policy1",
			Name:       "latest1",
			Severity:   storage.Severity_LOW_SEVERITY,
			Categories: []string{"Image Assurance", "Privileges Capabilities"},
			Fields: &storage.PolicyFields{
				ImageName: &storage.ImageNamePolicy{
					Tag: "latest1",
				},
			},
		},
		{
			Id:         "policy2",
			Name:       "latest2",
			Severity:   storage.Severity_LOW_SEVERITY,
			Categories: []string{"Image Assurance", "Privileges Capabilities"},
			Fields: &storage.PolicyFields{
				ImageName: &storage.ImageNamePolicy{
					Tag: "latest2",
				},
			},
		},
		{
			Id:         "policy3",
			Name:       "latest3",
			Severity:   storage.Severity_LOW_SEVERITY,
			Categories: []string{"Image Assurance", "Privileges Capabilities"},
			Fields: &storage.PolicyFields{
				ImageName: &storage.ImageNamePolicy{
					Tag: "latest3",
				},
			},
		},
	}
}
