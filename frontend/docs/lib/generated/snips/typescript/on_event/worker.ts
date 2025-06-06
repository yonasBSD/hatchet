import { Snippet } from '@/lib/generated/snips/types';

const snippet: Snippet = {
  "language": "typescript ",
  "content": "import { hatchet } from '../hatchet-client';\nimport { lower, upper } from './workflow';\n\nasync function main() {\n  const worker = await hatchet.worker('on-event-worker', {\n    workflows: [lower, upper],\n  });\n\n  await worker.start();\n}\n\nif (require.main === module) {\n  main();\n}\n",
  "source": "out/typescript/on_event/worker.ts",
  "blocks": {},
  "highlights": {}
};

export default snippet;
