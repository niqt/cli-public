# CLI-DEMO DEPENDENCIES

This application allow to get the dependencies information the deps.dev using the API. 
This information are saved in an sqlite db adding also the retrieved score (when it's possible)
The web application allow to search the dependencies by name or score (also using an interval).

## RUN
You need to install Docker and docker-compose on your computer, after that execute:
* docker-compose build
* docker-compose up -d

If everything is ok, browse to localhost (the web application run to the canonical 80 port)

The first thing to do is click on the Load button to download the dependencies, after that it's
possible do to the CRUD operation on that (NOTE: at the moment there is not validation on the edit dialog).

## API DESCRIPTION

------------------------------------------------------------------------------------------

#### Create a dependency

<details>
 <summary><code>POST</code> <code><b>/api/v1/package/{name}/{version}/dependency</b></code> <code>Add a dependency to the package with version</code></summary>

##### Parameters

> | name    |  type     | data type | description                       |
> |---------|-----------|-----------|-----------------------------------|
> | name    |  required | string    | the package name                  |
> | version |  required | string    | the the version of the package    |

##### Body
<code>
{id: integer, name: string, version: string, score: float, packageId: number}
</code>

##### Responses

> | http code | content-type              | response                    |
> |-----------|---------------------------|-----------------------------|
> | `201`     | `application/json`        | `{data:[dependencyObject]}` |
> | `500`     | `application/json`        | `{error: "message error}`   |

##### dependencyObject
{id: integer, name: string, version: string, score: float, packageId: number}


</details>

------------------------------------------------------------------------------------------

#### Get the dependencies

<details>
 <summary><code>GET</code> <code><b>/api/v1/package/{name}/{version}/dependency</b></code> <code>Get the dependencies of a package by version and name</code></summary>

##### Parameters

> | name       | type     | data type | description                                                           |
> |------------|----------|-----------|-----------------------------------------------------------------------|
> | name       | required | string    | the package name                                                      |
> | version    | required | string    | the version of the package                                            |
> | searchName | optional | string    | search only the dependencies that have searchName as part of the name |
> | lower      | optional | number    | search for a min score                                                |
> | upper      | optional | number    | search for a max score                                                |


##### Responses

> | http code | content-type              | response                    |
> |-----------|---------------------------|-----------------------------|
> | `200`     | `application/json`        | `{data:[dependencyObject]}` |
> | `500`     | `application/json`        | `{error: "message error}`   |

##### dependencyObject
{id: integer, name: string, version: string, score: float, packageId: number}


</details>

#### Update Dependency

<details>
 <summary><code>PUT</code> <code><b>"/api/v1/package/{name}/{version}/dependency/{id}"</b></code> <code>Update the dependency</code></summary>

##### Parameters

> | name    | type     | data type | description                    |
> |---------|----------|-----------|--------------------------------|
> | name    | required | string    | the package name               |
> | version | required | string    | the version of the package     |
> | id      | required | number    | id of the dependency to update |

##### Body
<code>
{id: integer, name: string, version: string, score: float, packageId: number}
</code>

##### Responses

> | http code | content-type              | response                    |
> |-----------|---------------------------|-----------------------------|
> | `200`     | `application/json`        | `{data:[dependencyObject]}` |
> | `500`     | `application/json`        | `{error: "message error}`   |

##### dependencyObject
{id: integer, name: string, version: string, score: float, packageId: number}

</details>

------------------------------------------------------------------------------------------

#### Delete Dependency

<details>
 <summary><code>DELETE</code> <code><b>"/api/v1/package/{name}/{version}/dependency/{id}"</b></code> <code>Delete the dependency</code></summary>

##### Parameters

> | name    | type     | data type | description                    |
> |---------|----------|-----------|--------------------------------|
> | name    | required | string    | the package name               |
> | version | required | string    | the version of the package     |
> | id      | required | number    | id of the dependency to delete |



##### Responses

> | http code | content-type              | response                  |
> |-----------|---------------------------|---------------------------|
> | `200`     | `application/json`        | `{data:[]}`               |
> | `500`     | `application/json`        | `{error: "message error}` |


</details>

------------------------------------------------------------------------------------------
#### Download the Dependencies from deps.dev

<details>
 <summary><code>PUT</code> <code><b>"/api/v1/package/{name}/{version}"</b></code> <code>Get the dependencies from deps.dev</code></summary>

##### Parameters

> | name    | type     | data type | description                    |
> |---------|----------|-----------|--------------------------------|
> | name    | required | string    | the package name               |
> | version | required | string    | the version of the package     |


##### Responses

> | http code | content-type              | response                    |
> |-----------|---------------------------|-----------------------------|
> | `200`     | `application/json`        | `{data:[dependencyObject]}` |
> | `500`     | `application/json`        | `{error: "message error}`   |

##### dependencyObject
{id: integer, name: string, version: string, score: float, packageId: number}

</details>

------------------------------------------------------------------------------------------

#### Get Package

<details>
 <summary><code>GET</code> <code><b>"/api/v1/package/{name}/{version}"</b></code> <code>Get the package information</code></summary>

##### Parameters

> | name       | type     | data type | description                                                           |
> |------------|----------|-----------|-----------------------------------------------------------------------|
> | name       | required | string    | the package name                                                      |
> | version    | required | string    | the version of the package                                            |



##### Responses

> | http code | content-type              | response                  |
> |-----------|---------------------------|---------------------------|
> | `200`     | `application/json`        | `{data:[packageObject]}`  |
> | `500`     | `application/json`        | `{error: "message error}` |

##### packageObject
{id: integer, name: string, version: string}


</details>

------------------------------------------------------------------------------------------

## DATA MODEL

The SQLite database uses only two tables:

##### Package

> | Col name | Type    | Not null | Auto Increment |
> |----------|---------|----------|----------------|
> | id       | INTEGER | Y        | Y              |
> | name     | TEXT    |          |                |
> | version  | TEXT    |          |                |

##### Dependency

> | Col name  | Type    | Not null | Auto Increment |
> |-----------|---------|----------|----------------|
> | id        | INTEGER | Y        | Y              |
> | name      | TEXT    |          |                |
> | version   | TEXT    |          |                |
> | score     | NUMERIC |          |                |
> | packageId | Integer |          |                |