const EventEmitter = {
  events: {},
  dispatch: function (event: string, data: unknown): void {
    if (!this.events[event]) return;
    this.events[event].forEach((callback: (arg0: unknown) => unknown) =>
      callback(data)
    );
  },
  subscribe: function (event: string, callback: unknown): void {
    if (!this.events[event]) this.events[event] = [];
    this.events[event].push(callback);
  },
  unsubscribe: function (event: string): void {
    if (!this.events[event]) return;
    this.events[event].splice(event);
  },
};

export default EventEmitter;
