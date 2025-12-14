from setuptools import setup, find_packages

setup(
    name="etherply",
    version="0.1.0",
    description="EtherPly Python SDK for real-time synchronization",
    author="EtherPly",
    packages=find_packages(),
    install_requires=[
        "websockets>=10.0",
        "aiohttp>=3.8.0",
    ],
    python_requires=">=3.8",
)
