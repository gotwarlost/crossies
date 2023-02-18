(function (window) {
    class LayoutTxn {
        constructor(tiles, layout) {
            this._tiles = tiles;
            this._layout = layout;
        }

        apply() {
            this._tiles.syncToLayout(this._layout);
        }
    }

    class ChangePhraseTxn {
        constructor(phrase, inputNode) {
            this._phrase = phrase;
            this._inputNode = inputNode;
        }

        apply() {
            this._inputNode.value = this._phrase;
        }
    }

    class BatchTxn {
        constructor(txns) {
            this._txns = txns;
        }

        apply() {
            this._txns.forEach(function (t) {
                t.apply();
            });
        }
    }

    class TxnStack {
        constructor(observer, cleanTxn) {
            this._txns = [
                cleanTxn,
            ];
            this._current = 0;
            this._observer = observer;
        }

        apply(txn) {
            txn.apply();
            this._txns.splice(this._current + 1);
            this._txns.push(txn);
            this._current = this._txns.length - 1;
            this._observer(this);
        }

        canUndo() {
            return this._current > 0;
        }

        undo() {
            if (!this.canUndo()) {
                console.warn("attempt to undo when not possible");
                return;
            }
            this._current -= 1;
            const t = this._txns[this._current];
            t.apply();
            this._observer(this);
        }

        canRedo() {
            return this._current < this._txns.length - 1;
        }

        redo() {
            if (!this.canRedo()) {
                console.warn("attempt to redo when not possible");
                return;
            }
            this._current += 1;
            const t = this._txns[this._current];
            t.apply();
            this._observer(this);
        }
    }
    window.LayoutTxn = LayoutTxn;
    window.ChangePhraseTxn = ChangePhraseTxn;
    window.BatchTxn = BatchTxn;
    window.TxnStack = TxnStack;
})(window);

