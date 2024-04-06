import JPLError from './error';
import JPLExecutionError from './execution';

/** Wrap occuring errors in a JPLExecutionError if they are not already a JPLError */
function adaptErrors(cb) {
  try {
    return cb();
  } catch (err) {
    throw adaptError(err);
  }
}

export default adaptErrors;

/** Wrap occuring errors in a JPLExecutionError if they are not already a JPLError */
export async function adaptErrorsAsync(cb) {
  try {
    return await cb();
  } catch (err) {
    throw adaptError(err);
  }
}

/** Wrap error in a JPLExecutionError if it is not already a JPLError */
export function adaptError(error) {
  if (!JPLError.is(error)) return new JPLExecutionError(error.message);
  return error;
}
