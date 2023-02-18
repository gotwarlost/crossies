(function (exports, document, PlainDraggable) {

    function calculateBoundingBox(node) {
        const rect = node.getBoundingClientRect();
        return {
            left: rect.left + window.scrollX,
            right: rect.right + window.scrollX,
            top: rect.top + window.scrollY,
            bottom: rect.bottom + window.scrollY,
        };
    }

    // Coords is a "coordinate system" used to convert absolute positions
    // (i.e. positions relative to the document) to relative ones (i.e. relative
    // to a dom element) and vice-versa.
    class Coords {
        constructor(boundingNode) {
            this._bn = boundingNode;
        }

        get boundingNode() {
            return this._bn;
        }

        boundingBox() {
            return calculateBoundingBox(this._bn);
        }

        toAbsolute(pt) {
            const rect = this.boundingBox();
            return {x: pt.x + rect.left, y: pt.y + rect.top};
        }

        toRelative(pt) {
            const rect = this.boundingBox();
            return {x: pt.x - rect.left, y: pt.y - rect.top};
        }
    }

    function createTileState(value, position) {
        return {
            value: value,
            position: position
        };
    }

    // Layout is a set of tile states and an associated bounding box for which
    // the tile positions are valid.
    class Layout {
        constructor(boundingRect, tileStates) {
            this._bb = boundingRect;
            this._tileStates = tileStates || [];
        }

        get tileStates() {
            return this._tileStates;
        }

        get boundingRect() {
            return this._bb;
        }

        get state() {
            return {
                boundingRect: this.boundingRect,
                tileStates: this.tileStates
            };
        }
    }

    class TileAnimation {
        constructor(draggable, duration) {
            this._draggable = draggable;
            this._duration = duration;
        }

        _animate(animFn, endFn, interval) {
            animFn();
            let t = 0;
            const endTime = Date.now() + this._duration;
            t = setInterval(function () {
                if (Date.now() > endTime) {
                    clearInterval(t);
                    endFn(t);
                    return;
                }
                animFn();
            }, interval);
        }

        create() {
            const el = this._draggable.element;
            const endFn = function () {
                el.style.opacity = 1;
            };
            const iterations = [0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9];
            const interval = this._duration / iterations.length;
            let count = 0;
            const animFn = function () {
                const opacity = iterations[count];
                if (count < iterations.length - 1) {
                    count += 1;
                }
                el.style.opacity = opacity;
            };
            this._animate(animFn, endFn, interval);
        }

        move(finalPosition) {
            const draggable = this._draggable;
            const r = draggable.rect;
            const endFn = function () {
                draggable.setOptions({
                    left: finalPosition.x,
                    top: finalPosition.y
                });
            };
            const interval = 10;
            const iterations = this._duration / interval;
            const xDistance = (finalPosition.x - r.left) / iterations;
            const yDistance = (finalPosition.y - r.top) / iterations;
            let count = 0;
            const animFn = function () {
                count += 1;
                draggable.setOptions({
                    left: r.left + Math.floor(count * xDistance),
                    top: r.top + Math.floor(count * yDistance)
                });
            };
            draggable.position();
            this._animate(animFn, endFn, interval);
        }

        setValue(newValue) {
            const el = this._draggable.element;
            const iterations = [0.8, 0.6, 0.4, 0.2, 0, 0.2, 0.4, 0.6, 0.8, 1];
            const interval = duration / iterations.length;
            let count = 0;
            const animFn = function () {
                const opacity = iterations[count];
                if (count < iterations.length - 1) {
                    count += 1;
                }
                el.style.opacity = opacity;
                if (opacity === 0) {
                    el.replaceChild(document.createTextNode(newValue), el.firstChild);
                }
            };
            const endFn = function () {
                el.style.opacity = 1;
            };
            this._animate(animFn, endFn, interval);
        }

        destroy() {
            const draggable = this._draggable;
            const el = draggable.element;
            const endFn = function () {
                draggable.remove();
                el.parentNode.removeChild(el);
            };
            const iterations = [0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1, 0];
            const interval = this._duration / iterations.length;
            let count = 0;
            const animFn = function () {
                const opacity = iterations[count];
                if (count < iterations.length - 1) {
                    count += 1;
                }
                el.style.opacity = opacity;
            };
            this._animate(animFn, endFn, interval);
        }
    }

    // Tile encapsulates the DOM element and corresponding draggable
    // associated with a tile.
    class Tile {
        constructor(coords, tileState, observer) {
            const el = document.createElement('div');
            el.setAttribute('class', 'tile-base tile-computed');
            el.setAttribute('style', 'opacity: 0');
            el.appendChild(document.createTextNode(tileState.value));
            coords.boundingNode.appendChild(el);
            this._id = Tile.nextID();
            this._value = tileState.value;
            this._coords = coords;
            const abs = coords.toAbsolute(tileState.position);
            const that = this;
            this._draggable = new PlainDraggable(el, {
                left: abs.x,
                top: abs.y,
                onDragEnd: function () {
                    observer(that);
                }
            });
            new TileAnimation(this._draggable, Tile.animationInterval).create();
        }

        static nextID() {
            Tile._nextID += 1;
            return 't' + Tile._nextID;
        }

        get DOMElement() {
            return this._draggable.element;
        }

        get id() {
            return this._id;
        }

        get value() {
            return this._value;
        }

        set value(v) {
            if (this._value === v) {
                return;
            }
            this._value = v;
            new TileAnimation(this._draggable, Tile.animationInterval).setValue(v);
        }

        get position() {
            const rect = this._draggable.rect;
            return this._coords.toRelative({x: rect.left, y: rect.top});
        }

        set position(pt) {
            const abs = this._coords.toAbsolute(pt);
            new TileAnimation(this._draggable, Tile.animationInterval).move(abs);
        }

        get state() {
            return createTileState(this.value, this.position);
        }

        destroy() {
            new TileAnimation(this._draggable, Tile.animationInterval).destroy();
        }
    }

    Tile._nextID = 0;
    Tile.animationInterval = 100;

    // Tiles is collection of tile objects that can "sync" to a layout.
    // Tile objects are created and destroyed as necessary to match the
    // desired layout.
    class Tiles {
        constructor(coords) {
            this._coords = coords;
            this._tiles = [];
            this._moveObserver = null;
        }

        set moveObserver(o) {
            this._moveObserver = o;
        }

        get tiles() {
            return this._tiles;
        }

        syncToLayout(layout) {
            const used = {}; // track tiles that have been matched up to a want state
            const unused = {}; // track tiles that haven't matched up to a want state
            const allocatedIDs = {}; // track want states that have been allocated
            const that = this;

            this._tiles.forEach(function (t) {
                unused[t.id] = t;
            });
            const wantStates = {};
            for (let i = 0; i < layout.tileStates.length; i += 1) {
                wantStates['w' + i] = layout.tileStates[i];
            }
            const matchValue = function (haveTile, wantState) {
                return haveTile.value === wantState.value;
            };
            const matchPosition = function (haveTile, wantState) {
                const hPos = haveTile.position;
                const wPos = wantState.position;
                return hPos.x === wPos.x && hPos.y === wPos.y;
            };
            const matchBoth = function (haveTile, wantState) {
                return matchValue(haveTile, wantState) && matchPosition(haveTile, wantState);
            };
            const matchers = [matchBoth, matchPosition, matchValue];
            const markUsed = function (tileID, wantID) {
                used[tileID] = unused[tileID];
                delete unused[tileID];
                allocatedIDs[wantID] = true;
            };
            matchers.forEach(function (matcher) {
                Object.keys(wantStates).forEach(function (wantId) {
                    const wantState = wantStates[wantId];
                    Object.keys(unused).forEach(function (haveId) {
                        if (allocatedIDs[wantId]) {
                            return;
                        }
                        const haveTile = unused[haveId];
                        if (!matcher(haveTile, wantState)) {
                            return;
                        }
                        if (haveTile.value !== wantState.value) {
                            haveTile.value = wantState.value;
                        }
                        const hPos = haveTile.position;
                        const wPos = wantState.position;
                        if (hPos.x !== wPos.x || hPos.y !== wPos.y) {
                            haveTile.position = wPos;
                        }
                        markUsed(haveId, wantId);
                    });
                });
            });
            const unmatchedStates = [];
            Object.keys(wantStates).forEach(function (id) {
                if (allocatedIDs[id]) {
                    return;
                }
                unmatchedStates.push(wantStates[id]);
            });

            const tiles = [];
            for (const k in used) {
                tiles.push(used[k]);
            }

            for (const k in unused) {
                const t = unused[k];
                t.destroy();
            }
            const coords = this._coords;
            const tileObserver = function (tile) {
                if (that._moveObserver) {
                    that._moveObserver(tile);
                }
            };
            unmatchedStates.forEach(function (state) {
                const t = new Tile(coords, state, tileObserver);
                tiles.push(t);
            });
            this._tiles = tiles;

            if (this._tiles.length !== layout.tileStates.length) {
                console.error("Layout", layout, "Tiles", tiles);
                throw "syncToLayout: invariant violated";
            }
        }

        currentLayout() {
            const states = [];
            this.tiles.forEach(function (tile) {
                states.push(tile.state);
            });
            return new Layout(this._coords.boundingBox(), states);
        }
    }

    exports.Coords = Coords;
    exports.createTileState = createTileState;
    exports.Tile = Tile;
    exports.Tiles = Tiles;
    exports.Layout = Layout;

})(window, document, window.PlainDraggable);
