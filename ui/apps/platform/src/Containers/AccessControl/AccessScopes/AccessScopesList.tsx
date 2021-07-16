import React, { ReactElement, useState } from 'react';
import {
    Alert,
    AlertVariant,
    Badge,
    Button,
    Modal,
    ModalVariant,
    Title,
    Toolbar,
    ToolbarContent,
    ToolbarGroup,
    ToolbarItem,
} from '@patternfly/react-core';
import { TableComposable, Tbody, Td, Thead, Th, Tr } from '@patternfly/react-table';

import { AccessScope, getIsDefaultAccessScopeId } from 'services/AccessScopesService';
import { Role } from 'services/RolesService';

import { AccessControlEntityLink, RolesLink } from '../AccessControlLinks';

const entityType = 'ACCESS_SCOPE';

export type AccessScopesListProps = {
    accessScopes: AccessScope[];
    roles: Role[];
    handleCreate: () => void;
    handleDelete: (id: string) => Promise<void>;
};

function AccessScopesList({
    accessScopes,
    roles,
    handleCreate,
    handleDelete,
}: AccessScopesListProps): ReactElement {
    const [idDeleting, setIdDeleting] = useState('');
    const [nameConfirmingDelete, setNameConfirmingDelete] = useState<string | null>(null);
    const [alertDelete, setAlertDelete] = useState<ReactElement | null>(null);

    function onClickDelete(id: string) {
        setIdDeleting(id);
        setNameConfirmingDelete(
            accessScopes.find((accessScope) => accessScope.id === id)?.name ?? ''
        );
    }

    function onConfirmDelete() {
        setNameConfirmingDelete(null);
        setAlertDelete(null);
        handleDelete(idDeleting)
            .catch((error) => {
                setAlertDelete(
                    <Alert
                        title="Delete access scope failed"
                        variant={AlertVariant.danger}
                        isInline
                    >
                        {error.message}
                    </Alert>
                );
            })
            .finally(() => {
                setIdDeleting('');
            });
    }

    function onCancelDelete() {
        setNameConfirmingDelete(null);
        setIdDeleting('');
    }

    return (
        <>
            <Toolbar inset={{ default: 'insetNone' }}>
                <ToolbarContent>
                    <ToolbarGroup spaceItems={{ default: 'spaceItemsMd' }}>
                        <ToolbarItem>
                            <Title headingLevel="h2">Access scopes</Title>
                        </ToolbarItem>
                        <ToolbarItem>
                            <Badge isRead>{accessScopes.length}</Badge>
                        </ToolbarItem>
                    </ToolbarGroup>
                    <ToolbarItem alignment={{ default: 'alignRight' }}>
                        <Button variant="primary" onClick={handleCreate} isSmall>
                            Add access scope
                        </Button>
                    </ToolbarItem>
                </ToolbarContent>
            </Toolbar>
            {alertDelete}
            <TableComposable variant="compact">
                <Thead>
                    <Tr>
                        <Th width={20}>Name</Th>
                        <Th width={30}>Description</Th>
                        <Th width={40}>Roles</Th>
                        <Th width={10} aria-label="Row actions" />
                    </Tr>
                </Thead>
                <Tbody>
                    <Tr className="pf-u-background-color-200">
                        <Td dataLabel="Name">Unrestricted</Td>
                        <Td dataLabel="Description">Access to all clusters and namespaces</Td>
                        <Td dataLabel="Roles">
                            <RolesLink
                                roles={roles.filter(
                                    ({ accessScopeId }) => accessScopeId.length === 0
                                )}
                                entityType={entityType}
                                entityId=""
                            />
                        </Td>
                        <Td />
                    </Tr>
                    {accessScopes.map(({ id, name, description }) => (
                        <Tr key={id}>
                            <Td dataLabel="Name">
                                <AccessControlEntityLink
                                    entityType={entityType}
                                    entityId={id}
                                    entityName={name}
                                />
                            </Td>
                            <Td dataLabel="Description">{description}</Td>
                            <Td dataLabel="Roles">
                                <RolesLink
                                    roles={roles.filter(
                                        ({ accessScopeId }) => accessScopeId === id
                                    )}
                                    entityType={entityType}
                                    entityId={id}
                                />
                            </Td>
                            <Td
                                actions={{
                                    disable:
                                        idDeleting === id ||
                                        getIsDefaultAccessScopeId(id) ||
                                        roles.some(({ accessScopeId }) => accessScopeId === id),
                                    items: [
                                        {
                                            title: 'Delete access scope',
                                            onClick: () => onClickDelete(id),
                                        },
                                    ],
                                }}
                                className="pf-u-text-align-right"
                            />
                        </Tr>
                    ))}
                </Tbody>
            </TableComposable>
            <Modal
                variant={ModalVariant.small}
                title="Permanently delete access scope?"
                isOpen={typeof nameConfirmingDelete === 'string'}
                onClose={onCancelDelete}
                actions={[
                    <Button key="confirm" variant="danger" onClick={onConfirmDelete}>
                        Delete
                    </Button>,
                    <Button key="cancel" variant="link" onClick={onCancelDelete}>
                        Cancel
                    </Button>,
                ]}
            >
                {nameConfirmingDelete ? (
                    <div>
                        Access scope name: <strong>{nameConfirmingDelete}</strong>
                    </div>
                ) : (
                    ''
                )}
            </Modal>
        </>
    );
}

export default AccessScopesList;
