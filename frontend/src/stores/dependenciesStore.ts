import { mande } from 'mande'
import { defineStore } from 'pinia'
import type {DependenciesForAppDto, DependencyDto} from "@/model/dependeciesDto.ts";
import type {PackageDto} from "@/model/packageDto.ts";


// Define interface for API error
interface ApiError {
  message: string
  code?: string
  status?: number
}

const packageName = 'github.com/cli/cli/v2';
const packageVersion = 'v2.66.1';

const encodedName = encodeURIComponent(packageName);
const encodedVersion = encodeURIComponent(packageVersion);

const url = `http://localhost:8080/api/v1/package/${encodedName}/${encodedVersion}/dependency`;
export const useDependencies = defineStore('dependencies', {
  state: () => ({
    dependencies: null as  DependencyDto[] | null,
    isLoading: false,
    error: null as ApiError | null,
  }),

  actions: {
    async getDependencies(searchText: string, min: string, max: string): Promise<ApiError | DependenciesForAppDto> {
      this.isLoading = true
      this.error = null

      try {
        const fullPath = url + `?searchName=${searchText}&lower=${min}&upper=${max}`;
        const api = mande(fullPath)
        this.isLoading = false
        const result = await api.get() as DependenciesForAppDto;
        if (result.data) {
          this.dependencies = result.data;
        } else {
          this.dependencies = [];
        }
        return result;
      } catch (err) {
        const error = err as ApiError
        this.error = error
        return error
      } finally {
        this.isLoading = false
      }
    },
    async updateDependency(dep: DependencyDto): Promise<ApiError | DependenciesForAppDto> {
      this.isLoading = true
      this.error = null

      try {

        const fullPath = `http://localhost:8080/api/v1/package/${encodedName}/${encodedVersion}/dependency/${dep.id}`;
        const api = mande(fullPath)
        this.isLoading = false
        return await api.put(dep) as DependenciesForAppDto;
      } catch (err) {
        const error = err as ApiError
        this.error = error
        return error
      } finally {
        this.isLoading = false
      }
    },
    async createDependency(dep: DependencyDto): Promise<ApiError | DependenciesForAppDto> {
      this.isLoading = true
      this.error = null

      try {
        const fullPath = `http://localhost:8080/api/v1/package/${encodedName}/${encodedVersion}/dependency`;
        console.log("FULL PATH: " + fullPath)
        const api = mande(fullPath)
        this.isLoading = false
        return await api.post(dep) as DependenciesForAppDto;
      } catch (err) {
        const error = err as ApiError
        this.error = error
        return error
      } finally {
        this.isLoading = false
      }
    },
    async getPackage(name: string, version: string): Promise<ApiError | PackageDto> {
      this.isLoading = true
      this.error = null
      const encodedName = encodeURIComponent(name);
      const encodedVersion = encodeURIComponent(version);
      try {
        const fullPath = `http://localhost:8080/api/v1/package/${encodedName}/${encodedVersion}`;
        const api = mande(fullPath)
        this.isLoading = false
        return await api.get() as PackageDto;
      } catch (err) {
        const error = err as ApiError
        this.error = error
        return error
      } finally {
        this.isLoading = false
      }
    },
    async downloadDependency(): Promise<ApiError | DependenciesForAppDto> {
      this.isLoading = true
      this.error = null

      try {
        const fullPath = `http://localhost:8080/api/v1/package/${encodedName}/${encodedVersion}`;
        const api = mande(fullPath)
        this.isLoading = false
        return await api.get() as DependenciesForAppDto;
      } catch (err) {
        const error = err as ApiError
        this.error = error
        return error
      } finally {
        this.isLoading = false
      }
    },
    async deleteDependency(id: number): Promise<ApiError | any> {
      this.isLoading = true
      this.error = null

      try {
        const fullPath = `http://localhost:8080/api/v1/package/${encodedName}/${encodedVersion}/dependency/${id}`;
        const api = mande(fullPath)
        this.isLoading = false
        return await api.delete() as any;
      } catch (err) {
        const error = err as ApiError
        this.error = error
        return error
      } finally {
        this.isLoading = false
      }
    },
  },
})

// Export type for components to use
export type DependenciesStore = ReturnType<typeof useDependencies>
