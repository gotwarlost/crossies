(function (window) {
    /*
        a snapshot is a plain JSON-serializable object that holds enough state to
        recreate the current word with its positions and bounding box.
     */
    function createSnapshot(phrase, layout) {
        return {
            phrase: phrase,
            layoutState: layout.state
        }
    }

    class Store {
        constructor() {
            this._currentSnapshot = null;
            this._subscribers = [];
        }

        get currentSnapshot() {
            return this._currentSnapshot;
        }

        set currentSnapshot(v) {
            this._currentSnapshot = v;
            this._subscribers.forEach(function (observer) {
                observer();
            })
        }
        addSubscriber(fn) {
            this._subscribers.push(fn);
        }
    }

    window.createSnapshot = createSnapshot;
    window.Store = Store;

})(window);
