/** Apply the specified defaults to the config */
function applyDefaults(config, defaults, ...strictPaths) {
  const iter = (c, d, s) => {
    if (Array.isArray(c)) return c;

    if (typeof c === 'object') {
      if (c === null || Array.isArray(d) || typeof d !== 'object' || d === null) return c;
      const result = { ...c };
      const isStrict = s.some((entry) => entry.length === 0);
      Object.keys({ ...c, ...d }).forEach((key) => {
        if (isStrict) {
          if (!Object.hasOwn(result, key)) result[key] = d[key];
        } else {
          result[key] = iter(
            result[key],
            d[key],
            s.filter((entry) => entry[0] === key).map((entry) => entry.slice(1)),
          );
        }
        if (result[key] === undefined) delete result[key];
      });
      return result;
    }

    return c === undefined ? d : c;
  };

  return iter(
    config,
    defaults,
    strictPaths.map((path) => (typeof path === 'string' ? path.split('.') : path)),
  );
}

export default applyDefaults;
