import { Snippet } from '@/lib/generated/snips/types';

const snippet: Snippet = {
  "language": "python",
  "content": "from examples.sticky_workers.worker import sticky_workflow\nfrom hatchet_sdk import TriggerWorkflowOptions\n\nsticky_workflow.run(\n    options=TriggerWorkflowOptions(additional_metadata={\"hello\": \"moon\"}),\n)\n",
  "source": "out/python/sticky_workers/event.py",
  "blocks": {},
  "highlights": {}
};

export default snippet;
