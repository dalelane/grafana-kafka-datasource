import type { DataSourceJsonData } from '@grafana/data';
import type { DataQuery } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  topicName: string;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  topicName: 'TOPIC_NAME',
};

export interface MyDataSourceOptions extends DataSourceJsonData {
  bootstrapservers: string;
  clientid: string;
  groupid: string;
  authtype: string;
  username: string;
  usetls: boolean;
}

export interface MySecureDataSourceOptions {
  password: string;
}
