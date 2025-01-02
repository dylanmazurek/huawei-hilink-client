"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const yargs_1 = __importDefault(require("yargs"));
const helpers_1 = require("yargs/helpers");
const ListSMS_1 = require("./src/ListSMS");
const MobileData_1 = require("./src/MobileData");
const Signal_1 = require("./src/Signal");
const startSession_1 = require("./src/startSession");
const huawei = require('./jslib/public');
function delay(ms) {
    return __awaiter(this, void 0, void 0, function* () {
        return new Promise(res => setTimeout(res, ms));
    });
}
// @ts-ignore
(0, yargs_1.default)((0, helpers_1.hideBin)(process.argv))
    .command('sendSMS', 'send SMS to contact or group of contacts', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('phone', {
        describe: 'phones with ; as separator ',
        type: "string",
    }).positional('message', {
        describe: 'text message ',
        type: "string",
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        if (!argv.phone) {
            throw new Error('Phone number is not defined');
        }
        yield (0, ListSMS_1.sendMessage)(sessionData, argv.phone, argv.message || '');
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('contacts', 'get contact list with the latest sms messages', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('page', {
        describe: 'sms page',
        default: 1
    }).positional('exportFile', {
        describe: 'export to file',
        default: './contacts.list'
    }).positional('exportFormat', {
        describe: 'export format (xml, json, none)',
        default: 'none'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        switch (argv.exportFormat) {
            case 'json': {
                break;
            }
            case 'xml': {
                break;
            }
            case 'none': {
                break;
            }
            default: {
                throw new Error(`export Format ${argv.exportFile} does not supported: supported only: xml,json,none`);
            }
        }
        yield (0, ListSMS_1.getSMSContacts)(sessionData, argv.page, argv.exportFile, argv.exportFormat);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('messages', 'get all messages from InBox', (yargs) => {
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    }).positional('deleteAfter', {
        describe: 'delete all messages after reading ',
        default: false
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    })
        .positional('exportFile', {
        describe: 'export to file',
        default: './inbox.list'
    }).positional('exportFormat', {
        describe: 'export format (xml, json, none)',
        default: 'none'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        switch (argv.exportFormat) {
            case 'json': {
                break;
            }
            case 'none': {
                break;
            }
            default: {
                throw new Error(`export Format ${argv.exportFile} does not supported: supported only: json,none`);
            }
        }
        yield (0, ListSMS_1.getInBoxSMS)(sessionData, argv.deleteAfter, argv.exportFile, argv.exportFormat);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('contactPages', 'contact list pages', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('exportFile', {
        describe: 'export to file',
        default: './contactsCount.list'
    }).positional('exportFormat', {
        describe: 'export format (xml, json, none)',
        default: 'none'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        switch (argv.exportFormat) {
            case 'json': {
                break;
            }
            case 'xml': {
                break;
            }
            case 'none': {
                break;
            }
            default: {
                throw new Error(`export Format ${argv.exportFile} does not supported: supported only: xml,json,none`);
            }
        }
        yield (0, ListSMS_1.getSMSPages)(sessionData, argv.exportFile, argv.exportFormat);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('sms', 'get contact SMS list', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1',
    }).positional('phone', {
        describe: 'contact phone number',
        type: 'string'
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('page', {
        describe: 'sms page',
        default: 1
    }).positional('exportFile', {
        describe: 'export to file',
        default: './sms.list'
    }).positional('exportFormat', {
        describe: 'export format (xml, json, none)',
        default: 'none'
    }).positional('deleteAfter', {
        describe: 'delete all messages after reading ',
        default: false
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        if (!argv.phone) {
            throw new Error('phone is not defined');
        }
        switch (argv.exportFormat) {
            case 'json': {
                break;
            }
            case 'xml': {
                break;
            }
            case 'none': {
                break;
            }
            default: {
                throw new Error(`export Format ${argv.exportFile} does not supported: supported only: xml,json,none`);
            }
        }
        yield (0, ListSMS_1.getSMSByUsers)(sessionData, argv.phone, argv.page, argv.exportFile, argv.exportFormat, argv.deleteAfter);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('pages', 'count of sms pages', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('phone', {
        describe: 'contact phone number',
        type: 'string'
    }).positional('exportFile', {
        describe: 'export to file',
        default: './smsCount.list'
    }).positional('exportFormat', {
        describe: 'export format (xml, json, none)',
        default: 'none'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        if (!argv.phone) {
            throw new Error('phone is not defined');
        }
        switch (argv.exportFormat) {
            case 'json': {
                break;
            }
            case 'xml': {
                break;
            }
            case 'none': {
                break;
            }
            default: {
                throw new Error(`export Format ${argv.exportFile} does not supported: supported only: xml,json,none`);
            }
        }
        yield (0, ListSMS_1.getContactSMSPages)(sessionData, argv.phone, argv.exportFile, argv.exportFormat);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('deleteSMS', 'delete sms by smsId', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    })
        .positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('messageId', {
        describe: 'messageId or index',
        type: 'string'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        if (!argv.messageId) {
            throw new Error('messageId is not defined');
        }
        yield (0, ListSMS_1.deleteMessage)(sessionData, argv.messageId);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('mobileData', 'Enable/Disable or Reconnect Mobile Data', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    }).positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('mode', {
        describe: 'change mobile data to on,off or reconnect',
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        switch (argv.mode) {
            case 'reconnect': {
                yield (0, MobileData_1.reconnect)(sessionData);
                return;
            }
            case 'on': {
                break;
            }
            case 'off': {
                break;
            }
            default: {
                throw new Error('Does not support Mode: ' + argv.mode + '. Supported only on,off,reconnect');
            }
        }
        yield (0, MobileData_1.controlMobileData)(sessionData, argv.mode);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('monitoring', 'current Monitoring status', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    }).positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('exportFile', {
        describe: 'export to file',
        default: './monitoring.log'
    }).positional('exportFormat', {
        describe: 'export format (xml, json, none)',
        default: 'none'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        const sessionData = yield (0, startSession_1.startSession)(argv.url);
        switch (argv.exportFormat) {
            case 'json': {
                break;
            }
            case 'xml': {
                break;
            }
            case 'none': {
                break;
            }
            default: {
                throw new Error(`export Format ${argv.exportFile} does not supported: supported only: xml,json,none`);
            }
        }
        yield (0, MobileData_1.status)(sessionData, argv.exportFile, argv.exportFormat);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('signalInfo', 'current device signal status', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    }).positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('turn', {
        alias: 't',
        describe: 'Request Signal info until interrupted',
        default: false,
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        //const sessionData = await startSession(argv.url);
        const sessionData = {
            url: argv.url,
            TokInfo: huawei.publicSession.token2,
            SesInfo: huawei.publicSession.session,
        };
        if (argv.turn) {
            while (true) {
                yield (0, Signal_1.getSignalInfo)(sessionData);
                yield delay(2500);
            }
        }
        else {
            yield (0, Signal_1.getSignalInfo)(sessionData);
        }
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
})).command('changeLteBand', 'change LTE band', (yargs) => {
    // @ts-ignore
    return yargs
        .positional('url', {
        describe: 'huawei host',
        default: '192.168.8.1'
    }).positional('username', {
        describe: 'huawei username',
        default: 'admin'
    })
        .positional('password', {
        describe: 'huawei password',
        type: 'string'
    }).positional('band', {
        alias: 'band',
        describe: 'desirable LTE band number, separated by + char (example 1+3+20).If you want to use every supported bands, write \'AUTO\'.", "AUTO"',
        default: 'AUTO',
        type: 'string'
    });
}, (argv) => __awaiter(void 0, void 0, void 0, function* () {
    if (!argv.band) {
        throw new Error('Band is empty');
    }
    huawei.publicKey.rsapadingtype = argv.rsapadingtype || "1";
    yield (0, startSession_1.login)(argv.url, argv.username, argv.password);
    try {
        yield (0, Signal_1.lteBand)(yield (0, startSession_1.startSession)(argv.url), argv.band);
    }
    finally {
        yield (0, startSession_1.logout)(argv.url);
    }
}))
    .option('rsapadingtype', {
    type: 'string',
    description: 'rsapadingtype, to check your run in web-console: MUI.LoginStateController.rsapadingtype',
    default: '1'
})
    .parse();
//# sourceMappingURL=index.js.map