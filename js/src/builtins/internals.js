import { nativeFunction } from '../library';

export const sortEntries = nativeFunction(async (runtime, input) => [
  await runtime.alterValue(input, (value) =>
    [...value].sort((a, b) =>
      runtime.compareArrays(runtime.unwrapValue(a)[0], runtime.unwrapValue(b)[0]),
    ),
  ),
]);
