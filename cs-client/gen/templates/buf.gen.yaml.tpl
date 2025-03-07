version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: {{ .PackagePrefix }}
clean: true
plugins:
  - local: protoc-gen-cosmos-csharp
    out: {{ .OutDir }}
  - remote: buf.build/protocolbuffers/csharp
    out: {{ .OutDir }}
    opt: file_extension=.pb.cs,base_namespace=
  - remote: buf.build/grpc/csharp
    out: {{ .OutDir }}
    opt: no_server,file_suffix=Grpc.pb.cs,base_namespace=
inputs:
  - directory: {{ .InDir }}
