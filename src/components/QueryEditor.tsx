import React, { ChangeEvent } from 'react';
import { InlineField, Input, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onTopicNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, topicName: event.target.value });
  };

  const { topicName } = query;

  return (
    <Stack>
      <InlineField label="Topic name" labelWidth={20} tooltip="topic to consume from">
        <Input onChange={onTopicNameChange} onBlur={onRunQuery} value={topicName || ''} type="string" />
      </InlineField>
    </Stack>
  );
}
