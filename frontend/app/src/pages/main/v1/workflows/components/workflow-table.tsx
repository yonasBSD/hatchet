import { useEffect, useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Workflow, queries } from '@/lib/api';
import invariant from 'tiny-invariant';
import { TenantContextType } from '@/lib/outlet';
import { Link, useOutletContext, useSearchParams } from 'react-router-dom';
import { DataTable } from '@/components/v1/molecules/data-table/data-table.tsx';
import { columns } from './workflow-columns';
import { Loading } from '@/components/v1/ui/loading.tsx';
import { Button } from '@/components/v1/ui/button';
import { ArrowPathIcon } from '@heroicons/react/24/outline';
import {
  PaginationState,
  SortingState,
  VisibilityState,
} from '@tanstack/react-table';
import { BiCard, BiTable } from 'react-icons/bi';
import RelativeDate from '@/components/v1/molecules/relative-date';
import { Badge } from '@/components/v1/ui/badge';
import { IntroDocsEmptyState } from '@/pages/onboarding/intro-docs-empty-state';

export function WorkflowTable() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { tenant } = useOutletContext<TenantContextType>();
  invariant(tenant);

  const [sorting, setSorting] = useState<SortingState>([
    {
      id: 'lastRun',
      desc: true,
    },
  ]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rotate, setRotate] = useState(false);

  const [cardToggle, setCardToggle] = useState(true);

  const [pagination, setPagination] = useState<PaginationState>(() => {
    const pageIndex = Number(searchParams.get('pageIndex')) || 0;
    const pageSize = Number(searchParams.get('pageSize')) || 50;

    return { pageIndex, pageSize };
  });
  const [pageSize, setPageSize] = useState<number>(
    Number(searchParams.get('pageSize')) || 50,
  );

  const [pageIndex, setPageIndex] = useState<number>(
    Number(searchParams.get('pageIndex')) || 0,
  );

  const listWorkflowQuery = useQuery({
    ...queries.workflows.list(tenant.metadata.id, {
      limit: pagination.pageSize,
      offset: pagination.pageIndex * pageSize,
    }),
    refetchInterval: 5000,
  });

  useEffect(() => {
    const newSearchParams = new URLSearchParams(searchParams);
    newSearchParams.set('pageIndex', pagination.pageIndex.toString());
    newSearchParams.set('pageSize', pagination.pageSize.toString());
    setSearchParams(newSearchParams);
  }, [pagination, setSearchParams, searchParams]);

  const data = useMemo(() => {
    const data = listWorkflowQuery.data?.rows || [];
    setPageIndex(listWorkflowQuery.data?.pagination?.num_pages || 0);

    return data;
  }, [
    listWorkflowQuery.data?.rows,
    listWorkflowQuery.data?.pagination?.num_pages,
  ]);

  if (listWorkflowQuery.isLoading) {
    return <Loading />;
  }

  const card: React.FC<{ data: Workflow }> = ({ data }) => (
    <div
      key={data.metadata?.id}
      className="border overflow-hidden shadow rounded-lg"
    >
      <div className="px-4 py-5 sm:p-6">
        <div className="flex flex-row justify-between items-center">
          <h3 className="text-lg leading-6 font-medium text-foreground">
            <Link to={`/v1/tasks/${data.metadata?.id}`}>{data.name}</Link>
          </h3>
          {data.isPaused ? (
            <Badge variant="inProgress">Paused</Badge>
          ) : (
            <Badge variant="successful">Active</Badge>
          )}
        </div>
        <p className="mt-1 max-w-2xl text-sm text-gray-700 dark:text-gray-300">
          {/* Last run{' '}
          {data.lastRunAt ? <RelativeDate date={data.lastRunAt} /> : 'never'}
          <br /> */}
          Created at <RelativeDate date={data.metadata?.createdAt} />
        </p>
      </div>
      <div className="px-4 py-4 sm:px-6">
        <div className="text-sm text-background-secondary">
          <Link to={`/v1/tasks/${data.metadata?.id}`}>
            <Button>View Workflow</Button>
          </Link>
        </div>
      </div>
    </div>
  );

  const actions = [
    <Button
      key="card-toggle"
      className="h-8 px-2 lg:px-3"
      size="sm"
      onClick={() => {
        setCardToggle((t) => !t);
      }}
      variant={'outline'}
      aria-label="Toggle card/table view"
    >
      {!cardToggle ? (
        <BiCard className={`h-4 w-4 `} />
      ) : (
        <BiTable className={`h-4 w-4 `} />
      )}
    </Button>,
    <Button
      key="refresh"
      className="h-8 px-2 lg:px-3"
      size="sm"
      onClick={() => {
        listWorkflowQuery.refetch();
        setRotate(!rotate);
      }}
      variant={'outline'}
      aria-label="Refresh events list"
    >
      <ArrowPathIcon
        className={`h-4 w-4 transition-transform ${rotate ? 'rotate-180' : ''}`}
      />
    </Button>,
  ];

  return (
    <DataTable
      columns={columns}
      data={data}
      filters={[]}
      emptyState={
        <IntroDocsEmptyState
          link="/home/your-first-task"
          title="No Registered Workflows"
          linkPreambleText="To learn more about workflows function in Hatchet,"
          linkText="check out our documentation."
        />
      }
      columnVisibility={columnVisibility}
      setColumnVisibility={setColumnVisibility}
      pagination={pagination}
      setPagination={setPagination}
      onSetPageSize={setPageSize}
      showSelectedRows={false}
      pageCount={pageIndex}
      sorting={sorting}
      setSorting={setSorting}
      isLoading={listWorkflowQuery.isLoading}
      manualSorting={false}
      actions={actions}
      manualFiltering={false}
      card={
        cardToggle
          ? {
              component: card,
            }
          : undefined
      }
    />
  );
}
