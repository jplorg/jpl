import JPLFatalError from './errors/fatal';

/** Runtime signal for managing execution lifecycle */
class RuntimeSignal {
  constructor(parent) {
    this._parent = parent;
    this._exited = false;
    this._subscriptions = {};
    this._nextSubscriptionKey = 0;
  }

  /** Return the signal's parent */
  get parent() {
    return this._parent;
  }

  /** Check whether the current runtime area has been requested to be exited */
  get exited() {
    this._exited ||= this.parent?.exited ?? false;
    return this._exited;
  }

  /** Check that the current runtime area has not been requested to be exited and throw a fatal error otherwise */
  checkHealth() {
    if (this.exited) throw new JPLFatalError('execution aborted');
  }

  /**
   * Request the current runtime area to be exited.
   * This also involves all child areas (introduced using `RuntimeSignal.next`).
   */
  exit = () => {
    if (this.exited) return;
    this._exited = true;
    Object.values(this._subscriptions).forEach((subscription) => {
      subscription();
    });
    this._subscriptions = {};
  };

  /**
   * Subscribe for when the current runtime area is requested to be exited.
   * This also involves all parent areas.
   * The function returns an unsubscription hook which must be called when completed in order to prevent memory leaks.
   * If the area has already been exited when subscribing, the callback is called immediately after the current event cycle.
   */
  subscribe = (cb) => {
    if (this.exited) {
      setTimeout(cb);
      return () => {};
    }
    const key = this._nextSubscriptionKey;
    this._nextSubscriptionKey += 1;
    this._subscriptions[key] = cb;
    const unsubscribeParent = (this.parent?.subscribe ?? (() => () => {}))(cb);
    const unsubscribe = () => {
      delete this._subscriptions[key];
      unsubscribeParent();
    };
    return unsubscribe;
  };

  /** Inherit the next child area for the current one */
  next = () => new RuntimeSignal(this);
}

export default RuntimeSignal;
