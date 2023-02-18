(function (window, Layout) {
    function shuffle(letters) {
        const count = letters.length;
        for (let i = 0; i < count; i += 1) {
            const r = Math.floor(Math.random() * count);
            const tmp = letters[i];
            letters[i] = letters[r];
            letters[r] = tmp;
        }
    }

    function randomGridLayout(letters, gridWidth, gridHeight) {
        const count = letters.length;
        shuffle(letters);
        const numCells = gridWidth * gridHeight;
        let ratio = Math.sqrt(count / numCells);
        let w = Math.floor(ratio * gridWidth);
        let h = Math.floor(ratio * gridHeight);
        if (w * h < count) {
            if (gridWidth > gridHeight) {
                w += 1;
            } else {
                h += 1;
            }
        }
        if (w * h < count) {
            if (gridWidth > gridHeight) {
                h += 1;
            } else {
                w += 1;
            }
        }
        const startX = Math.floor((gridWidth - w) / 2);
        const startY = Math.floor((gridHeight - h) / 2);
        const ret = [];
        for (let row = 0; row < h; row += 1) {
            for (let col = 0; col < w; col += 1) {
                const index = row * w + col;
                if (index < count) {
                    ret.push({value: letters[index], row: startY + row, col: startX + col});
                }
            }
        }
        return ret;
    }

    function randomLayout(phrase, coords) {
        const letters = phrase.split("").filter(function (ch) {
            return !(ch === ' ' || ch === '\t' || ch === '\n');
        });
        const cellEdge = 120;
        const realCellEdge = 80;
        const rect = coords.boundingBox();
        const gridWidth = Math.floor((rect.right - rect.left) / cellEdge);
        const gridHeight = Math.floor((rect.bottom - rect.top) / cellEdge);
        const states = [];

        const positions = randomGridLayout(letters, gridWidth, gridHeight);
        if (positions.length !== letters.length) {
            console.error('POSITIONS', positions, 'LETTERS', letters);
            throw "layout invariant violated!";
        }
        positions.forEach(function (pos) {
            const r1 = Math.floor(Math.random() * (cellEdge - realCellEdge));
            const r2 = Math.floor(Math.random() * (cellEdge - realCellEdge));
            const pt = {x: pos.col * cellEdge + r1, y: pos.row * cellEdge + r2};
            states.push(createTileState(pos.value, pt));
        });
        return new Layout(coords.boundingBox(), states);
    }

    window.randomLayout = randomLayout;
})(window, window.Layout);
