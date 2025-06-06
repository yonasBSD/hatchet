import { Snippet } from '@/lib/generated/snips/types';

const snippet: Snippet = {
  "language": "python",
  "content": "import time\n\nfrom examples.cancellation.worker import cancellation_workflow, hatchet\n\nid = cancellation_workflow.run_no_wait()\n\ntime.sleep(5)\n\nhatchet.runs.cancel(id.workflow_run_id)\n",
  "source": "out/python/cancellation/trigger.py",
  "blocks": {},
  "highlights": {}
};

export default snippet;
