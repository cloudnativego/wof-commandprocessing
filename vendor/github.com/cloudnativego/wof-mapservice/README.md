[![wercker status](https://app.wercker.com/status/bd75f3e11ac23cd6c10fa8a745da492b/m "wercker status")](https://app.wercker.com/project/bykey/bd75f3e11ac23cd6c10fa8a745da492b)

# World of FluxCraft - Map Service

This service exposes a RESTful API for the storage and retrieval of maps that can be used for individual games of World of FluxCraft.

API for this service is as follows:

| Resource | Method | Description |
|---|---|---|
| /api/maps | GET | Retrieves a list of all maps in the system, with full detail |
| /api/maps/**(id)** | GET | Retrieves the map with the supplied **ID** |
| /api/maps/**(id)** | PUT | Creates or Updates the map with the supplied **ID** | 
