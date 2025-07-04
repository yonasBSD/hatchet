# coding: utf-8

"""
    Hatchet API

    The Hatchet API

    The version of the OpenAPI document: 1.0.0
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


from __future__ import annotations

import json
from enum import Enum

from typing_extensions import Self


class TenantUIVersion(str, Enum):
    """
    TenantUIVersion
    """

    """
    allowed enum values
    """
    V0 = "V0"
    V1 = "V1"

    @classmethod
    def from_json(cls, json_str: str) -> Self:
        """Create an instance of TenantUIVersion from a JSON string"""
        return cls(json.loads(json_str))
