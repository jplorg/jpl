import JPLFatalError from './errors/fatal';

/** Runtime signal for managing execution lifecycle */
export default class JPLRuntimeSignal {
  #parent;
  #exited;
  #subscriptions: { [key: string]: () => void };
  #nextSubscriptionKey;

  constructor(parent?: JPLRuntimeSignal) {
    this.#parent = parent;
    this.#exited = false;
    this.#subscriptions = {};
    this.#nextSubscriptionKey = 0;
  }

  /** Return the signal's parent */
  get parent() {
    return this.#parent;
  }

  /** Check whether the current runtime area has been requested to be exited */
  get exited(): boolean {
    this.#exited ||= this.parent?.exited ?? false;
    return this.#exited;
  }

  /** Check that the current runtime area has not been requested to be exited and throw a fatal error otherwise */
  checkHealth(): void {
    if (this.exited) throw new JPLFatalError('execution aborted');
  }

  /**
   * Request the current runtime area to be exited.
   * This also involves all child areas (introduced using `JPLRuntimeSignal.next`).
   */
  exit = (): void => {
    if (this.exited) return;
    this.#exited = true;
    Object.values(this.#subscriptions).forEach((subscription) => {
      subscription();
    });
    this.#subscriptions = {};
  };

  /**
   * Subscribe for when the current runtime area is requested to be exited.
   * This also involves all parent areas.
   * The function returns an unsubscription hook which must be called when completed in order to prevent memory leaks.
   * If the area has already been exited when subscribing, the callback is called immediately after the current event cycle.
   */
  subscribe = (cb: () => void): (() => void) => {
    if (this.exited) {
      setTimeout(cb);
      return () => {};
    }
    const key = this.#nextSubscriptionKey;
    this.#nextSubscriptionKey += 1;
    this.#subscriptions[key] = cb;
    const unsubscribeParent = (this.parent?.subscribe ?? (() => () => {}))(cb);
    const unsubscribe = () => {
      delete this.#subscriptions[key];
      unsubscribeParent();
    };
    return unsubscribe;
  };

  /** Inherit the next child area for the current one */
  next = () => new JPLRuntimeSignal(this);
}
