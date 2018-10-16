package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stackrox/rox/cmd/roxdetect/report"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/images/utils"
	"golang.org/x/net/context"
)

const tokenEnv = "STACKROX_TOKEN"

var (
	version  = "development"
	central  = flag.String("central", "localhost:8443", "Host and port endpoint where Central is located.")
	image    = flag.String("image", "", "the image name and reference (e.g. nginx:latest or nginx@sha256:...)")
	severity = flag.String("severity", "critical", "Exit with a non-zero status if any violated policies meet or exceed this severity.\nAllowed values are "+levelText()+".\nIgnored if \"-json\" is set.")

	versionFlag = flag.Bool("version", false, `Prints the version "`+version+`" and exits.`)
	json        = flag.Bool("json", false, "Output policy results as json.")
)

func main() {
	// Parse the input flags.
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(2)
	}

	if err := mainCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "roxdetect: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainCmd() error {
	if err := checkLevel(*severity); err != nil {
		return err
	}

	// Read token from ENV.
	token, exists := os.LookupEnv(tokenEnv)
	if !exists {
		return errors.New("the STACKROX_TOKEN environment variable must be set to a token generated by stackrox for a Remote host")
	}

	// Get the violated policies for the input data.
	violatedPolicies, err := getViolatedPolicies(token)
	if err != nil {
		return err
	}

	// If json mode was given, print results (as json) and immediately return.
	if *json {
		return report.JSON(os.Stdout, violatedPolicies)
	}

	// Print results in human readable mode.
	if err = report.Pretty(os.Stdout, violatedPolicies); err != nil {
		return err
	}

	// Exit with a status of 1 if any of the violated policies were of a
	// exceeded our severity threshold level.
	for _, policy := range violatedPolicies {
		if exceedsThreshold(*severity, policy.Severity) {
			return errors.New("policy severity exceeds threshold")
		}
	}

	return nil
}

// Fetch the alerts for the inputs and convert them to a list of Policies that are violated.
func getViolatedPolicies(token string) ([]*v1.Policy, error) {
	alerts, err := getAlerts(token)
	if err != nil {
		return nil, err
	}

	var policies []*v1.Policy
	for _, alert := range alerts {
		policies = append(policies, alert.GetPolicy())
	}
	return policies, nil
}

// Get the alerts for the command line inputs.
func getAlerts(token string) ([]*v1.Alert, error) {
	// Attempt to construct the request first since it is the cheapest op.
	image, err := buildRequest()
	if err != nil {
		return nil, err
	}

	// Create the connection to the central detection service.
	conn, err := clientconn.UnauthenticatedGRPCConnection(*central)
	if err != nil {
		return nil, err
	}
	service := v1.NewDetectionServiceClient(conn)

	// Build context with token header.
	md := metautils.NiceMD{}
	md = md.Add("authorization", "Bearer "+token)
	ctx := md.ToOutgoing(context.Background())

	// Call detection and return the returned alerts.
	response, err := service.DetectBuildTime(ctx, image)
	if err != nil {
		return nil, err
	}
	return response.GetAlerts(), nil
}

// Use inputs to generate an image name for request.
func buildRequest() (*v1.Image, error) {
	if *image == "" {
		return nil, fmt.Errorf("image name must be set")
	}
	img, err := utils.GenerateImageFromStringWithError(*image)
	if err != nil {
		return nil, fmt.Errorf("could not parse image '%s': %s", *image, err)
	}
	return img, nil
}

// levelNames return an ordered list of string names that are derived from API
// severity levels. Additionally, the names "any" and "none" are added to this
// list as additional options.
func levelNames() []string {
	names := make([]string, len(v1.Severity_value)+2)
	names[0] = "any"
	names[len(names)-1] = "none"
	for key, value := range v1.Severity_value {
		name := strings.ToLower(strings.TrimSuffix(key, "_SEVERITY"))
		names[value+1] = name
	}
	return names
}

// levelText returns a textual description of the severity levels.
// Example: `"none", "unset", ..."critical", or "any"`
func levelText() string {
	names := levelNames()
	result := ""
	for index, name := range names {
		switch index {
		case 0:
			result += fmt.Sprintf(`"%s"`, name)
		case len(names) - 1:
			result += fmt.Sprintf(`, or "%s"`, name)
		default:
			result += fmt.Sprintf(`, "%s"`, name)
		}
	}
	return result
}

// checkLevel checks if the given level is allowed. An error is returned if an
// unknown level name is given.
func checkLevel(level string) error {
	names := levelNames()
	for _, name := range names {
		if level == name {
			return nil
		}
	}
	return fmt.Errorf("unknown severity level")
}

// exceedsThreshold checks if the given severity meets or exceeds the named
// threshold level.
func exceedsThreshold(threshold string, severity v1.Severity) bool {
	switch threshold {
	case "any":
		return true
	case "none":
		return false
	default:
		key := strings.ToUpper(threshold) + "_SEVERITY"
		return severity >= v1.Severity(v1.Severity_value[key])
	}
}
