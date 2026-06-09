import JPLError from './error';
import JPLExecutionError from './execution';

/** Wrap occuring errors in a JPLExecutionError if they are not already a JPLError */
export default function adaptErrors<T>(cb: () => T): T {
  try {
    return cb();
  } catch (err) {
    throw adaptError(err);
  }
}

/** Wrap occuring errors in a JPLExecutionError if they are not already a JPLError */
export async function adaptErrorsAsync<T>(cb: () => Promise<T>): Promise<T> {
  try {
    return await cb();
  } catch (err) {
    throw adaptError(err);
  }
}

/** Wrap error in a JPLExecutionError if it is not already a JPLError */
export function adaptError(error: unknown): unknown {
  if (!JPLError.is(error))
    return new JPLExecutionError((error as Error)?.message);
  return error;
}
