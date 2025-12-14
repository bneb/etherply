from typing import TypedDict, Any, Dict, Optional, Literal

class EtherPlyConfig(TypedDict, total=False):
    workspace_id: str
    token: str
    user_id: Optional[str]
    host: Optional[str]
    secure: Optional[bool]
    auto_reconnect: Optional[bool]
    max_reconnect_attempts: Optional[int]
    reconnect_base_delay: Optional[int]

class OperationPayload(TypedDict):
    key: str
    value: Any
    timestamp: int

class EtherPlyMessage(TypedDict):
    type: Literal['init', 'op', 'error']
    payload: Optional[OperationPayload]
    data: Optional[Dict[str, Any]]
    message: Optional[str]
