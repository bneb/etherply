class EtherPlyError(Exception):
    """Base exception for all EtherPly errors."""
    pass

class ConfigurationError(EtherPlyError):
    """Raised when configuration is invalid."""
    pass

class ConnectionError(EtherPlyError):
    """Raised when connection fails."""
    pass

class MessageError(EtherPlyError):
    """Raised when a message is invalid."""
    pass
