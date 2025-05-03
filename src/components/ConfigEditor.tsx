import React, { ChangeEvent } from 'react';
import { FieldSet, InlineField, InlineFieldRow, Input, Select, Checkbox } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { MyDataSourceOptions, MySecureDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions, MySecureDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

  const onBootstrapServersChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        bootstrapservers: event.target.value,
      },
    });
  };

  const onClientIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        clientid: event.target.value,
      },
    });
  };

  const onGroupIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        groupid: event.target.value,
      },
    });
  };

  const onAuthTypeChange = (selected: SelectableValue<string>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        authtype: selected.value || 'none',
      },
    });
  };

  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        username: event.target.value,
      },
    });
  };

  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        password: event.target.value,
      },
    });
  };

  const onTlsEnabledChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        usetls: event.target.checked,
      },
    });
  };

  const { bootstrapservers, clientid, groupid, authtype, username, usetls } = options.jsonData;

  const AUTH_TYPES = [
    { label: 'None', value: 'none' },
    { label: 'SCRAM-SHA-256', value: 'SCRAM-SHA-256' },
    { label: 'SCRAM-SHA-512', value: 'SCRAM-SHA-512' },
  ];

  return (
    <FieldSet label="Connection">
      <InlineFieldRow>
        <InlineField label="Bootstrap servers" labelWidth={30} tooltip="host:port">
          <Input
            width={50}
            required
            value={bootstrapservers}
            autoComplete="off"
            placeholder="localhost:9092"
            onChange={onBootstrapServersChange}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Client Id" labelWidth={30}>
          <Input
            width={25}
            required
            value={clientid}
            autoComplete="off"
            placeholder="grafana"
            onChange={onClientIdChange}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Consumer Group Id" labelWidth={30}>
          <Input
            width={25}
            required
            value={groupid}
            autoComplete="off"
            placeholder="grafana"
            onChange={onGroupIdChange}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Authentication type" labelWidth={20}>
          <Select
            required
            width={25}
            isMulti={false}
            backspaceRemovesValue={false}
            isClearable={false}
            defaultValue={AUTH_TYPES[0]}
            value={authtype}
            options={AUTH_TYPES}
            onChange={onAuthTypeChange}
          />
        </InlineField>
        <InlineField label="username" labelWidth={15}>
          <Input
            width={25}
            value={username}
            autoComplete="off"
            disabled={authtype === 'none'}
            onChange={onUsernameChange}
          />
        </InlineField>
        <InlineField label="password" labelWidth={15}>
          <Input
            width={25}
            type="password"
            value={options.secureJsonData?.password}
            autoComplete="off"
            disabled={authtype === 'none'}
            onChange={onPasswordChange}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Use TLS">
          <Checkbox value={usetls} checked={usetls} onChange={onTlsEnabledChange} />
        </InlineField>
      </InlineFieldRow>
    </FieldSet>
  );
}
