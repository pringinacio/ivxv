# IVXV Internet voting framework
"""
Common schema validator classes.
"""

from schematics.exceptions import ValidationError
from schematics.models import Model
from schematics.types import (IntType, ListType, ModelType, PolyModelType,
                              StringType, URLType)

from .fields import CertificateType


# FIXME: constrainments on schemas passed to wrapper aren't enforced.
# fields can  be missing from the conf, type constrainments are ignored,
# without raising any errors (e.g missing storage.conf.ca,
# invalid URL in qualification.*.url)
def protocol_cfg(mapping, **kwargs):
    """
    Return an alternative protocol configuration type.

    Alternative protocol configurations allow only one from a selection of
    alternatives to be used. The model looks like this::

        protocol: <protocol name>
        conf: <configuration for the chosen protocol>

    schematics has PolyModelType which allows "conf" to be different models,
    but it must select the model to use based solely on the contents of "conf",
    which cannot be done with this construction. So create a new wrapper model
    for each protocol configuration model which also includes the "protocol"
    field and let PolyModelType choose from those.
    """

    models = []
    for protocol, model in mapping.items():
        wrapper = type(
            f"{model.__name__}Wrapper", (Model, ), {
                "protocol": StringType(required=True, choices=[protocol]),
                "conf": ModelType(model, required=True),
                "ordertimeout": IntType(required=False, min_value=1)
            })
        mapping[protocol] = wrapper
        models.append(wrapper)

    def _claim(_, data):
        return mapping.get(data.get("protocol"))

    return PolyModelType(models, claim_function=_claim, **kwargs)


class DummySchema(Model):
    """Validating schema for Dummy container config."""
    trusted = ListType(StringType)


class OCSPSchema(Model):
    """Validating schema for OCSP config."""
    url = URLType(required=True)
    responders = ListType(CertificateType)
    retry = IntType(default=0, min_value=0)
    maxSkew = IntType(default=300, min_value=0)  # 300 milliseconds
    maxAge = IntType(default=1, min_value=0)  # 1 minute
    maxAge = IntType(default=1, min_value=0)  # 1 minute


class OCSPSchemaNoURL(Model):
    """Validating schema for OCSP config."""
    responders = ListType(CertificateType)
    maxSkew = IntType(default=300, min_value=0)  # 300 milliseconds
    maxAge = IntType(default=1, min_value=0)  # 1 minute


class TSPSchema(Model):
    """Validating schema for timestamp protocol config."""
    url = URLType(required=True)
    signers = ListType(CertificateType, required=True)
    delaytime = IntType(required=True, min_value=0)
    retry = IntType(default=0, min_value=0)
    maxSkew = IntType(default=2, min_value=0)  # 2 seconds
    maxAge = IntType(default=1, min_value=0)  # 1 minute


class TSPSchemaNoURL(Model):
    """Validating schema for timestamp protocol config."""
    signers = ListType(CertificateType, required=True)
    delaytime = IntType(required=True, min_value=0)
    maxSkew = IntType(default=2, min_value=0)  # 2 seconds
    maxAge = IntType(default=1, min_value=0)  # 1 minute


class BDocSchema(Model):
    """Validating schema for BDoc config."""
    bdocsize = IntType(required=True, min_value=1)
    filesize = IntType(required=True, min_value=1)
    roots = ListType(CertificateType, required=True)
    intermediates = ListType(CertificateType)
    profile = StringType(choices=['BES', 'TM', 'TS'], required=True)
    ocsp = ModelType(OCSPSchemaNoURL)
    tsp = ModelType(TSPSchemaNoURL)
    tsdelaytime = IntType(default=0, min_value=0)

    def validate_tsp(self, data, value):
        """Check that tsp exists if profile is TS."""
        try:
            if (data['profile'] == 'TS' and not data['tsp']):
                raise ValidationError('TS profile requires a tsp block')
        except KeyError:
            pass  # error in data structure is catched later
        return value


class ContainerSchema(Model):
    """Validating schema for signed container config."""
    bdoc = ModelType(BDocSchema)
    dummy = ModelType(DummySchema)
