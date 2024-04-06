function createSingleton(generatorFn) {
  let result;
  return () => {
    if (result == null) result = generatorFn();
    return result;
  };
}

export default createSingleton;
