<Project Sdk="Microsoft.NET.Sdk">
	<PropertyGroup>
		<PackageId>{{ .Name }}</PackageId>
		<Version>0.1,0</Version>
		<Authors>lxgr</Authors>
		<TargetFramework>netstandard2.1</TargetFramework>
		<Nullable>enable</Nullable>
		<PackageDescription>A C# client library for the {{ .ShortName }} blockchain</PackageDescription>
		<RepositoryUrl>{{ .URL }}</RepositoryUrl>
	</PropertyGroup>

	<PropertyGroup>
		<PackageReadmeFile>README.md</PackageReadmeFile>
	</PropertyGroup>

	<ItemGroup>
		<None Include="README.md" Pack="true" PackagePath="\"/>
	</ItemGroup>

	<PropertyGroup>
		<PackageLicenseExpression>MIT</PackageLicenseExpression>
	</PropertyGroup>

	<ItemGroup>
		<PackageReference Include="Cosmcs" Version="0.7.0-rc2"/>
	</ItemGroup>
</Project>
