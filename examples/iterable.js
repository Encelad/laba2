'use strict';

const iterable = {
    [Symbol.iterator]() {
        let i = 0;
        const iterator = {
            next() {
                return {
                    value: i++,
                    done: i > 3
                };
            }
        };
        return iterator;
    }
};


module.exports = {iterable};