import { nativeFunction } from '../library';

export const sortEntries = nativeFunction(async (runtime, input) => [
  await runtime.alterValue(input, (value) =>
    [...(value as unknown[])].sort((a, b) =>
      runtime.compareArrays(
        (runtime.unwrapValue(a) as unknown[])[0],
        (runtime.unwrapValue(b) as unknown[])[0],
      ),
    ),
  ),
]);
