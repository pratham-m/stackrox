import React from 'react';
import PropTypes from 'prop-types';

import Loader from 'Components/Loader';
import { labelClassName } from 'constants/form.constants';
import useFeatureFlagEnabled from 'hooks/useFeatureFlagEnabled';
import { knownBackendFlags } from 'utils/featureFlags';

import ClusterSummary from './Components/ClusterSummary';
import StaticConfigurationSection from './StaticConfigurationSection';
import DynamicConfigurationSection from './DynamicConfigurationSection';
import ClusterLabelsTable from './ClusterLabelsTable';

function ClusterEditForm({
    centralEnv,
    centralVersion,
    selectedCluster,
    handleChange,
    handleChangeLabels,
    isLoading,
}) {
    const hasScopedAccessControl = useFeatureFlagEnabled(
        knownBackendFlags.ROX_SCOPED_ACCESS_CONTROL
    );
    // guard against missing health status in unconnected and new clusters
    const healthStatus = selectedCluster.healthStatus || {
        overallHealthStatus: 'UNINITIALIZED',
    };

    if (isLoading) {
        return <Loader />;
    }

    return (
        <div className="bg-base-200 px-4 w-full">
            {/* @TODO, replace open prop with dynamic logic, based on clusterType */}
            {selectedCluster.id && (
                <ClusterSummary
                    healthStatus={healthStatus}
                    status={selectedCluster.status}
                    centralVersion={centralVersion}
                    currentDatetime={new Date()}
                    clusterId={selectedCluster.id}
                />
            )}
            <form
                className="grid grid-columns-1 md:grid-columns-2 grid-gap-4 xl:grid-gap-6 mb-4 w-full"
                data-testid="cluster-form"
            >
                <StaticConfigurationSection
                    centralEnv={centralEnv}
                    handleChange={handleChange}
                    selectedCluster={selectedCluster}
                />
                <div>
                    <DynamicConfigurationSection
                        dynamicConfig={selectedCluster.dynamicConfig}
                        helmConfig={selectedCluster.helmConfig}
                        handleChange={handleChange}
                    />
                    {hasScopedAccessControl && (
                        <div className="pt-4">
                            <label htmlFor="labels" className={labelClassName}>
                                Cluster labels
                            </label>
                            <ClusterLabelsTable
                                labels={selectedCluster?.labels ?? {}}
                                handleChangeLabels={handleChangeLabels}
                                hasAction
                            />
                        </div>
                    )}
                </div>
            </form>
        </div>
    );
}

ClusterEditForm.propTypes = {
    centralEnv: PropTypes.shape({
        kernelSupportAvailable: PropTypes.bool,
        successfullyFetched: PropTypes.bool,
    }).isRequired,
    centralVersion: PropTypes.string.isRequired,
    selectedCluster: PropTypes.shape({
        id: PropTypes.string,
        name: PropTypes.string,
        type: PropTypes.string,
        mainImage: PropTypes.string,
        centralApiEndpoint: PropTypes.string,
        collectionMethod: PropTypes.string,
        collectorImage: PropTypes.string,
        admissionController: PropTypes.bool,
        admissionControllerUpdates: PropTypes.bool,
        tolerationsConfig: PropTypes.shape({
            disabled: PropTypes.bool,
        }),
        status: PropTypes.shape({
            sensorVersion: PropTypes.string,
            providerMetadata: PropTypes.shape({
                region: PropTypes.string,
            }),
            orchestratorMetadata: PropTypes.shape({
                version: PropTypes.string,
                buildDate: PropTypes.string,
            }),
            upgradeStatus: PropTypes.shape({
                upgradability: PropTypes.string,
                upgradabilityStatusReason: PropTypes.string,
                mostRecentProcess: PropTypes.shape({
                    active: PropTypes.bool,
                    progress: PropTypes.shape({
                        upgradeState: PropTypes.string,
                        upgradeStatusDetail: PropTypes.string,
                    }),
                    type: PropTypes.string,
                }),
            }),
            certExpiryStatus: PropTypes.shape({
                sensorCertExpiry: PropTypes.string,
            }),
        }),
        dynamicConfig: PropTypes.shape({
            registryOverride: PropTypes.string,
            admissionControllerConfig: PropTypes.shape({
                enabled: PropTypes.bool,
                enforceOnUpdates: PropTypes.bool,
                timeoutSeconds: PropTypes.number,
                scanInline: PropTypes.bool,
                disableBypass: PropTypes.bool,
            }),
        }),
        helmConfig: PropTypes.shape({
            staticConfig: PropTypes.shape({}),
            dynamicConfig: PropTypes.shape({}),
        }),
        slimCollector: PropTypes.bool,
        healthStatus: PropTypes.shape({
            collectorHealthInfo: PropTypes.shape({
                version: PropTypes.string,
                totalDesiredPods: PropTypes.number,
                totalReadyPods: PropTypes.number,
                totalRegisteredNodes: PropTypes.number,
                statusErrors: PropTypes.arrayOf(PropTypes.string),
            }),
            sensorHealthStatus: PropTypes.string,
            collectorHealthStatus: PropTypes.string,
            overallHealthStatus: PropTypes.string,
            lastContact: PropTypes.string, // ISO 8601
            healthInfoComplete: PropTypes.bool,
        }),
        labels: PropTypes.shape({}),
    }).isRequired,
    handleChange: PropTypes.func.isRequired,
    handleChangeLabels: PropTypes.func.isRequired,
    isLoading: PropTypes.bool.isRequired,
};

export default ClusterEditForm;
