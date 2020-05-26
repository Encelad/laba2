'use strict';

const dir = (() => {
    let draw = "I'm in\n";
    return (arg) => {
        console.dir(arg);
        console.log(draw);
        draw = "Let me out!\n";
    }
})();

const log = arg => {
    console.log(arg);
}


module.exports = {
    console: {
        dir,
        log
    }
}
