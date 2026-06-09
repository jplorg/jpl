export default function createSingleton<T>(generatorFn: () => T): () => T {
  let result: T;
  return () => {
    if (result == null) result = generatorFn();
    return result;
  };
}
