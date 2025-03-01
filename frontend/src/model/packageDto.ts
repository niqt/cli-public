export interface Package {
  id: number;
  name: string;
  version: string;
}

export interface PackageDto {
  data?: Package;
  error?: string;
}
