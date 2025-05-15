import React from "react";
import { Divider, Input, Typography } from "antd";
import { Handle, Position } from "reactflow";

export const TimerNode: React.FC<{ data: any, id: string }> = ({ data, id }) => {
  return (
    <div
      style={{
        minWidth: 180,
        maxWidth: 280,
        background: "#fafafa",
        padding: `0px 5px`,
        borderRadius: 5,
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <Typography.Text strong>{data.description}</Typography.Text>
      </div>
      <Divider style={{ margin: "0" }} />
      <div
        style={{
          marginTop: 5,
          padding: "5px 0px",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <Typography.Text style={{ color: "#666" }}>
          {"Входные параметры"}
        </Typography.Text>
        {data.argIn.map((val: string) => {
          return (
            <Input
              addonBefore={val}
              size="small"
              style={{ padding: "2px 0px" }}
              value={data.input?.[val] || ""}
              onChange={e => data.onInputChange?.(id, val, e.target.value)}
            />
          );
        })}
        <Typography.Text style={{ color: "#666" }}>
          {"Выходные параметры"}
        </Typography.Text>
        {data.argOut.map((val: string) => {
          return (
            <Input
              addonBefore={val}
              size="small"
              style={{ padding: "2px 0px" }}
              value={data.output?.[val] || ""}
              onChange={e => data.onOutputChange?.(id, val, e.target.value)}
            />
          );
        })}
      </div>
      <Handle type="target" position={Position.Left} />
      <Handle type="source" position={Position.Right} />
    </div>
  );
};
