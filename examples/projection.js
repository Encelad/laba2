'use strict';

// Projection

const id = x => x;

const projection = meta => src => meta.reduce(
    (dest, [name, fn = id, field = name]) =>
        (dest[name] = fn(src[field]), dest), {}
);

// Display


const max = items => Math.max(...items);
const maxProp = key => items => max(items.map(x => x[key]));
const maxLength = maxProp('length');
const col = (name, data) => data.map(obj => obj[name].toString());

const render = meta => src => {
    const keys = meta.map(([name]) => name);
    const maxWidth = keys.map(key => maxLength(col(key, src)));
    const dest = src.map(obj => maxWidth.map(
        (width, i) => obj[keys[i]].toString().padEnd(width + 4)
    ));
    return dest.map(row => row.join('')).join('\n');
};

module.exports = {
    projection,
    render
}