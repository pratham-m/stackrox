import React, { useEffect, useState } from 'react';
import { Alert, Flex, FlexItem, Card, CardBody, Title } from '@patternfly/react-core';

import { fetchDeployment } from 'services/DeploymentsService';
import { portExposureLabels } from 'messages/common';
import ObjectDescriptionList from 'Components/ObjectDescriptionList';
import DeploymentOverview from './Deployment/DeploymentOverview';
import SecurityContext from './Deployment/SecurityContext';
import ContainerConfigurations from './Deployment/ContainerConfigurations';

type PortExposure = 'EXTERNAL' | 'NODE' | 'HOST' | 'INTERNAL' | 'UNSET';

type Port = {
    exposure: PortExposure;
    exposureInfos: {
        externalHostnames: string[];
        externalIps: string[];
        level: PortExposure;
        nodePort: number;
        serviceClusterIp: string;
        serviceId: string;
        serviceName: string;
        servicePort: number;
    }[];
    containerPort: number;
    exposedPort: number;
    name: string;
    protocol: string;
};

type FormattedPort = {
    exposure: string;
    exposureInfos: {
        externalHostnames: string[];
        externalIps: string[];
        level: string;
        nodePort: number;
        serviceClusterIp: string;
        serviceId: string;
        serviceName: string;
        servicePort: number;
    }[];
    containerPort: number;
    exposedPort: number;
    name: string;
    protocol: string;
};

export const formatDeploymentPorts = (ports: Port[] = []): FormattedPort[] => {
    const formattedPorts = [] as FormattedPort[];
    ports.forEach(({ exposure, exposureInfos, ...rest }) => {
        const formattedPort = { ...rest } as FormattedPort;
        formattedPort.exposure = portExposureLabels[exposure] || portExposureLabels.UNSET;
        formattedPort.exposureInfos = exposureInfos.map(({ level, ...restInfo }) => {
            return { ...restInfo, level: portExposureLabels[level] };
        });
        formattedPorts.push(formattedPort);
    });
    return formattedPorts;
};

const DeploymentDetails = ({ deployment }) => {
    // attempt to fetch related deployment to selected alert
    const [relatedDeployment, setRelatedDeployment] = useState(deployment);

    useEffect(() => {
        fetchDeployment(deployment.id).then(
            (dep) => setRelatedDeployment(dep),
            () => setRelatedDeployment(null)
        );
    }, [deployment.id, setRelatedDeployment]);

    const deploymentObj = relatedDeployment || deployment;

    return (
        <Flex className="pf-u-mt-md">
            {!relatedDeployment && (
                <Alert
                    variant="warning"
                    isInline
                    title="This data is a snapshot of a deployment that no longer exists."
                />
            )}
            <Flex flex={{ default: 'flex_1' }}>
                <Flex direction={{ default: 'column' }} flex={{ default: 'flex_1' }}>
                    <FlexItem>
                        <Title headingLevel="h3">Overview</Title>
                    </FlexItem>
                    <FlexItem>
                        <Card isFlat>
                            <CardBody>
                                <DeploymentOverview deployment={deploymentObj} />
                            </CardBody>
                        </Card>
                    </FlexItem>
                    <FlexItem>
                        <Title headingLevel="h3">Port Configuration</Title>
                    </FlexItem>
                    <FlexItem>
                        <Card isFlat>
                            <CardBody>
                                {deploymentObj?.ports?.length > 0
                                    ? formatDeploymentPorts(deploymentObj.ports).map((port) => (
                                          <ObjectDescriptionList data={port} />
                                      ))
                                    : 'None'}
                            </CardBody>
                        </Card>
                    </FlexItem>
                    <FlexItem>
                        <Title headingLevel="h3">Security Context</Title>
                    </FlexItem>
                    <FlexItem>
                        <SecurityContext deployment={relatedDeployment} />
                    </FlexItem>
                </Flex>
            </Flex>
            <Flex direction={{ default: 'column' }} flex={{ default: 'flex_1' }}>
                <FlexItem>
                    <Title headingLevel="h3">Container Configuration</Title>
                </FlexItem>
                <FlexItem>
                    <ContainerConfigurations deployment={relatedDeployment} />
                </FlexItem>
            </Flex>
        </Flex>
    );
};

export default DeploymentDetails;