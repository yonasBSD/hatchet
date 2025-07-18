import { Button } from '@/components/v1/ui/button';
import { Separator } from '@/components/v1/ui/separator';
import { useState } from 'react';
import { CreateInviteForm } from './components/create-invite-form';
import { useApiError } from '@/lib/hooks';
import { useMutation, useQuery } from '@tanstack/react-query';
import api, {
  CreateTenantInviteRequest,
  TenantInvite,
  UpdateTenantInviteRequest,
  UserChangePasswordRequest,
  queries,
} from '@/lib/api';
import { Dialog } from '@/components/v1/ui/dialog';
import { DataTable } from '@/components/v1/molecules/data-table/data-table';
import { columns } from './components/invites-columns';
import { columns as membersColumns } from './components/members-columns';
import { UpdateInviteForm } from './components/update-invite-form';
import { DeleteInviteForm } from './components/delete-invite-form';
import { ChangePasswordDialog } from './components/change-password-dialog';
import { AxiosError } from 'axios';
import useApiMeta from '@/pages/auth/hooks/use-api-meta';
import { useCurrentTenantId } from '@/hooks/use-tenant';

export default function Members() {
  const meta = useApiMeta();

  return (
    <div className="flex-grow h-full w-full">
      <div className="mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <h2 className="text-2xl font-bold leading-tight text-foreground">
          Members and Invites
        </h2>
        <Separator className="my-4" />
        <MembersList />
        {meta.data?.allowInvites && (
          <>
            <Separator className="my-4" />
            <InvitesList />
          </>
        )}
      </div>
    </div>
  );
}

function MembersList() {
  const { tenantId } = useCurrentTenantId();
  const [showChangePasswordDialog, setShowChangePasswordDialog] =
    useState(false);

  const listMembersQuery = useQuery({
    ...queries.members.list(tenantId),
  });

  return (
    <div>
      <h3 className="text-xl font-semibold leading-tight text-foreground">
        Members
      </h3>
      <Separator className="my-4" />
      <DataTable
        columns={membersColumns({
          onChangePasswordClick: () => {
            setShowChangePasswordDialog(true);
          },
        })}
        data={listMembersQuery.data?.rows || []}
        filters={[]}
        getRowId={(row) => row.metadata.id}
        isLoading={listMembersQuery.isLoading}
      />
      {showChangePasswordDialog && (
        <ChangePassword
          showChangePasswordDialog={showChangePasswordDialog}
          setShowChangePasswordDialog={setShowChangePasswordDialog}
          onSuccess={() => {}}
        />
      )}
    </div>
  );
}

function InvitesList() {
  const { tenantId } = useCurrentTenantId();
  const [showCreateInviteModal, setShowCreateInviteModal] = useState(false);
  const [updateInvite, setUpdateInvite] = useState<TenantInvite | null>(null);
  const [deleteInvite, setDeleteInvite] = useState<TenantInvite | null>(null);

  const listInvitesQuery = useQuery({
    ...queries.invites.list(tenantId),
  });

  const cols = columns({
    onEditClick: (row) => {
      setUpdateInvite(row);
    },
    onDeleteClick: (row) => {
      setDeleteInvite(row);
    },
  });

  return (
    <div>
      <div className="flex flex-row justify-between items-center">
        <h3 className="text-xl font-semibold leading-tight text-foreground">
          Invites
        </h3>
        <Button
          key="create-invite"
          onClick={() => setShowCreateInviteModal(true)}
        >
          Create Invite
        </Button>
      </div>
      <Separator className="my-4" />
      <DataTable
        isLoading={listInvitesQuery.isLoading}
        columns={cols}
        data={listInvitesQuery.data?.rows || []}
        filters={[]}
        getRowId={(row) => row.metadata.id}
      />
      {showCreateInviteModal && (
        <CreateInvite
          showCreateInviteModal={showCreateInviteModal}
          setShowCreateInviteModal={setShowCreateInviteModal}
          onSuccess={() => {
            setShowCreateInviteModal(false);
            listInvitesQuery.refetch();
          }}
        />
      )}
      {updateInvite && (
        <UpdateInvite
          tenantInvite={updateInvite}
          setShowTenantInvite={() => setUpdateInvite(null)}
          onSuccess={() => {
            setUpdateInvite(null);
            listInvitesQuery.refetch();
          }}
        />
      )}
      {deleteInvite && (
        <DeleteInvite
          tenantInvite={deleteInvite}
          setShowTenantInviteDelete={() => setDeleteInvite(null)}
          onSuccess={() => {
            setDeleteInvite(null);
            listInvitesQuery.refetch();
          }}
        />
      )}
    </div>
  );
}

function CreateInvite({
  showCreateInviteModal,
  setShowCreateInviteModal,
  onSuccess,
}: {
  showCreateInviteModal: boolean;
  setShowCreateInviteModal: (show: boolean) => void;
  onSuccess: () => void;
}) {
  const { tenantId } = useCurrentTenantId();
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const { handleApiError } = useApiError({
    setFieldErrors: setFieldErrors,
  });

  const createMutation = useMutation({
    mutationKey: ['tenant-invite:create', tenantId],
    mutationFn: async (data: CreateTenantInviteRequest) => {
      await api.tenantInviteCreate(tenantId, data);
    },
    onSuccess: onSuccess,
    onError: handleApiError,
  });

  return (
    <Dialog
      open={showCreateInviteModal}
      onOpenChange={setShowCreateInviteModal}
    >
      <CreateInviteForm
        isLoading={createMutation.isPending}
        onSubmit={createMutation.mutate}
        fieldErrors={fieldErrors}
      />
    </Dialog>
  );
}

function UpdateInvite({
  tenantInvite,
  setShowTenantInvite,
  onSuccess,
}: {
  tenantInvite: TenantInvite;
  setShowTenantInvite: (show: boolean) => void;
  onSuccess: () => void;
}) {
  const { tenantId } = useCurrentTenantId();

  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const { handleApiError } = useApiError({
    setFieldErrors: setFieldErrors,
  });

  const updateMutation = useMutation({
    mutationKey: ['tenant-invite:update', tenantId, tenantInvite],
    mutationFn: async (data: UpdateTenantInviteRequest) => {
      await api.tenantInviteUpdate(tenantId, tenantInvite.metadata.id, data);
    },
    onSuccess: onSuccess,
    onError: handleApiError,
  });

  return (
    <Dialog open={!!tenantInvite} onOpenChange={setShowTenantInvite}>
      <UpdateInviteForm
        invite={tenantInvite}
        isLoading={updateMutation.isPending}
        onSubmit={updateMutation.mutate}
        fieldErrors={fieldErrors}
      />
    </Dialog>
  );
}

function DeleteInvite({
  tenantInvite,
  setShowTenantInviteDelete,
  onSuccess,
}: {
  tenantInvite: TenantInvite;
  setShowTenantInviteDelete: (show: boolean) => void;
  onSuccess: () => void;
}) {
  const { tenantId } = useCurrentTenantId();
  const { handleApiError } = useApiError({});

  const deleteMutation = useMutation({
    mutationKey: ['tenant-invite:delete', tenantId, tenantInvite],
    mutationFn: async () => {
      await api.tenantInviteDelete(tenantId, tenantInvite.metadata.id);
    },
    onSuccess: onSuccess,
    onError: handleApiError,
  });

  return (
    <Dialog open={!!tenantInvite} onOpenChange={setShowTenantInviteDelete}>
      <DeleteInviteForm
        invite={tenantInvite}
        isLoading={deleteMutation.isPending}
        onSubmit={() => deleteMutation.mutate()}
        onCancel={() => setShowTenantInviteDelete(false)}
      />
    </Dialog>
  );
}

function ChangePassword({
  showChangePasswordDialog,
  setShowChangePasswordDialog,
  onSuccess,
}: {
  onSuccess: () => void;
  showChangePasswordDialog: boolean;
  setShowChangePasswordDialog: (show: boolean) => void;
}) {
  const { tenantId } = useCurrentTenantId();
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const { handleApiError } = useApiError({
    setFieldErrors: setFieldErrors,
  });

  const updatePasswordMutation = useMutation({
    mutationKey: ['user:update', tenantId],
    mutationFn: async (data: UserChangePasswordRequest) => {
      const res = await api.userUpdatePassword(data);
      return res.data;
    },
    onMutate: () => {
      setFieldErrors({});
    },
    onSuccess: () => {
      onSuccess();
      setShowChangePasswordDialog(false);
    },
    onError: (e: AxiosError<unknown, any>) => {
      return handleApiError(e);
    },
  });

  return (
    <Dialog
      open={showChangePasswordDialog}
      onOpenChange={setShowChangePasswordDialog}
    >
      <ChangePasswordDialog
        isLoading={updatePasswordMutation.isPending}
        onSubmit={(data) =>
          updatePasswordMutation.mutate({
            password: data.password,
            newPassword: data.newPassword,
          })
        }
        fieldErrors={fieldErrors}
      />
    </Dialog>
  );
}
