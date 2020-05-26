'use strict';

const {iterable} = require('./iterable.js');
const {console} = require('./console.js');

//what about iterable?

{
    for (const element of iterable) {
        console.dir({element});
    }

}

{
    const array = Array.from(iterable);
    console.log({array});
    console.log('-------------------------\n');
}

{
    const array = [...iterable];
    console.dir({array});
}

//end of iterable

console.log('\n\n\nLets talks about projection\n\n\n');

const {projection, render} = require('./projection.js');

//what about projection

// Dataset
const persons = [
    {name: 'Dasha', city: 'Kyiv', born: 2000},
    {name: 'Person1', city: 'Somewhere', born: 1994},
    {name: 'Tin', city: 'Odessa-mama', born: 1999},
    {name: 'Ya', city: 'Ungwar', born: 1596537},
];
// Metadata

const year = date => date.getFullYear();
const diff = y => year(new Date()) - year(new Date(y + ''));
const upper = s => s.toUpperCase();

const md = [
    ['name'],
    ['place', upper, 'city'],
    ['age', diff, 'born'],
];

const query = person => (
    person.name !== '' &&
    person.born < 2000 &&
    person.city === 'Odessa-mama'
);

const pf = projection(md);
const data = persons.filter(query).map(pf);

const renderer = render(md);
const res = renderer(data);
console.log('\n' + res + '\n');
