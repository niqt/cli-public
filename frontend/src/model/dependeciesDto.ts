export interface DependencyDto {
  id: number;
  name: string;
  version: string;
  score: number;
  packageId: number;
}

export interface DependenciesForAppDto {
  data?: DependencyDto[];
  error?: string;
}
