import { hatchet } from '../hatchet-client';

async function main() {
  const event = await hatchet.events.push('user:event', {
    Data: { Hello: 'World' },
  });
}

if (require.main === module) {
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      // eslint-disable-next-line no-console
      console.error('Error:', error);
      process.exit(1);
    });
}
