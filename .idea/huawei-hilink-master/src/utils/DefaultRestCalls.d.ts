import { AxiosPromise } from 'axios';
import { HTTPMethod, RestCalls } from './restCalls';
declare class DefaultRestCalls implements RestCalls {
    fetchData(url: string, method: HTTPMethod, headers?: any): Promise<string>;
    fetchDataRaw(url: string, method: HTTPMethod, headers?: any): Promise<AxiosPromise>;
    sendData(url: string, method: HTTPMethod, data: string, headers?: any): Promise<string>;
    sendDataRaw(url: string, method: HTTPMethod, data: string, headers?: any): Promise<AxiosPromise>;
}
export declare const restCalls: DefaultRestCalls;
export {};
//# sourceMappingURL=DefaultRestCalls.d.ts.map