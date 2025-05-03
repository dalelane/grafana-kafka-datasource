import {
  CoreApp,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  LiveChannelScope,
} from '@grafana/data';
import { DataSourceWithBackend, getGrafanaLiveSrv } from '@grafana/runtime';

import { Observable, merge } from 'rxjs';

import { MyQuery, MyDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    const observables = request.targets.map((query) => {
      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          path: query.topicName,
          data: {
            ...query,
          },
        },
      });
    });

    return merge(...observables);
  }

  getDefaultQuery(_: CoreApp): Partial<MyQuery> {
    return DEFAULT_QUERY;
  }

  filterQuery(query: MyQuery): boolean {
    if (
      query.hide ||
      !query.topicName ||
      query.topicName.trim() === '' ||
      query.topicName === DEFAULT_QUERY.topicName
    ) {
      return false;
    }
    return true;
  }
}
