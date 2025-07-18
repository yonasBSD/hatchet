/* eslint-disable no-console */
import { StickyStrategy } from '@hatchet/protoc/workflows';
import { hatchet } from '../hatchet-client';
import { child } from '../child_workflows/workflow';

// > Sticky Task
export const sticky = hatchet.task({
  name: 'sticky',
  retries: 3,
  sticky: StickyStrategy.SOFT,
  fn: async (_, ctx) => {
    // specify a child workflow to run on the same worker
    const result = await child.run(
      {
        N: 1,
      },
      { sticky: true }
    );

    return {
      result,
    };
  },
});
// !!
