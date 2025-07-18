import { Button } from '@/components/v1/ui/button';
import { useTenant } from '@/lib/atoms';
import { ArrowTopRightIcon } from '@radix-ui/react-icons';
import { Link } from 'react-router-dom';

export const WorkflowDefinitionLink = ({
  workflowId,
}: {
  workflowId: string;
}) => {
  const { tenant } = useTenant();

  return (
    <Link
      to={`/tenants/${tenant?.metadata.id}/workflows/${workflowId}`}
      target="_blank"
      rel="noreferrer"
    >
      <Button size={'sm'} className="px-2 py-2 gap-2" variant="outline">
        <ArrowTopRightIcon className="w-4 h-4" />
        Configuration
      </Button>
    </Link>
  );
};
