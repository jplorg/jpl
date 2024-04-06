/** Compare the specified strings based on their unicode code points */
function compareStrings(a, b) {
  const len = Math.min(a.length, b.length);
  let i = 0;
  // eslint-disable-next-line no-unused-vars
  for (const _ of a) {
    if (i >= len) {
      break;
    }
    const cp1 = a.codePointAt(i);
    const cp2 = b.codePointAt(i);
    const order = cp1 - cp2;
    if (order !== 0) {
      return order;
    }
    i += 1;
    if (cp1 > 0xffff) {
      i += 1;
    }
  }
  return a.length - b.length;
}

export default compareStrings;
