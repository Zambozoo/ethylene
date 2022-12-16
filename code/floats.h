#ifndef clox_floats_h
#define clox_floats_h

#include <stdint.h>

typedef union float64 {
    uint64_t u;
    double f;
} float64;

typedef union float32 {
    uint32_t u;
    float f;
    struct
    {
        uint32_t Mantissa : 23;
        uint32_t Exponent : 8;
        uint32_t Sign : 1;
    };
} float32;

typedef union float16 {
    uint16_t u;
    struct
    {
        uint32_t Mantissa : 10;
        uint32_t Exponent : 5;
        uint32_t Sign : 1;
    };
} float16;

typedef union float8 {
    uint8_t u;
    struct
    {
        uint32_t Mantissa : 3;
        uint32_t Exponent : 4;
        uint32_t Sign : 1;
    };
} float8;

static float16 floatToHalf(float32 f) {
    float16 o = { 0 };

    // Based on ISPC reference code (with minor modifications)
    if (f.Exponent == 0) // Signed zero/denormal (which will underflow)
        o.Exponent = 0;
    else if (f.Exponent == 255) // Inf or NaN (all exponent bits set)
    {
        o.Exponent = 31;
        o.Mantissa = f.Mantissa ? 0x200 : 0; // NaN->qNaN and Inf->Inf
    }
    else // Normalized number
    {
        // Exponent unbias the single, then bias the halfp
        int newexp = f.Exponent - 127 + 15;
        if (newexp >= 31) // Overflow, return signed infinity
            o.Exponent = 31;
        else if (newexp <= 0) // Underflow
        {
            if ((14 - newexp) <= 24) // Mantissa might be non-zero
            {
                uint32_t mant = f.Mantissa | 0x800000; // Hidden 1 bit
                o.Mantissa = mant >> (14 - newexp);
                if ((mant >> (13 - newexp)) & 1) // Check for rounding
                    o.u++; // Round, might overflow into exp bit, but this is OK
            }
        }
        else
        {
            o.Exponent = newexp;
            o.Mantissa = f.Mantissa >> 13;
            if (f.Mantissa & 0x1000) // Check for rounding
                o.u++; // Round, might overflow to inf, this is OK
        }
    }

    o.Sign = f.Sign;
    return o;
}

static float32 halfToFloat(float16 h) {
    static const float32 magic = { 113 << 23 };
    static const uint32_t shifted_exp = 0x7c00 << 13; // exponent mask after shift
    float32 o;

    o.u = (h.u & 0x7fff) << 13;     // exponent/mantissa bits
    uint32_t exp = shifted_exp & o.u;   // just the exponent
    o.u += (127 - 15) << 23;        // exponent adjust

    // handle exponent special cases
    if (exp == shifted_exp) // Inf/NaN?
        o.u += (128 - 16) << 23;    // extra exp adjust
    else if (exp == 0) // Zero/Denormal?
    {
        o.u += 1 << 23;             // extra exp adjust
        o.f -= magic.f;             // renormalize
    }

    o.u |= (h.u & 0x8000) << 16;    // sign bit
    return o;
}

static float8 floatToQuarter(float32 f) {
     float8 o = { 0 };

    // Based on ISPC reference code (with minor modifications)
    if (f.Exponent == 0) // Signed zero/denormal (which will underflow)
        o.Exponent = 0;
    else if (f.Exponent == 255) // Inf or NaN (all exponent bits set)
    {
        o.Exponent = 15;
        o.Mantissa = f.Mantissa ? 0b100 : 0; // NaN->qNaN and Inf->Inf
    }
    else // Normalized number
    {
        // Exponent unbias the single, then bias the halfp
        int newexp = f.Exponent - 127 + 2;
        if (newexp >= 15) // Overflow, return signed infinity
            o.Exponent = 15;
        else if (newexp <= 0) // Underflow
        {
            if ((1 - newexp) <= 24) // Mantissa might be non-zero
            {
                uint32_t mant = f.Mantissa | 0x800000; // Hidden 1 bit
                o.Mantissa = mant >> (1 - newexp);
                if ((mant >> (-newexp)) & 1) // Check for rounding
                    o.u++; // Round, might overflow into exp bit, but this is OK
            }
        }
        else
        {
            o.Exponent = newexp;
            o.Mantissa = f.Mantissa >> 21;
            if (f.Mantissa & 0x1000) // Check for rounding
                o.u++; // Round, might overflow to inf, this is OK
        }
    }

    o.Sign = f.Sign;
    return o;
}
static float32 quarterToFloat(float8 f) {
    static const float32 magic = { 113 << 23 };
    static const uint32_t shifted_exp = 0x7c00 << 13; // exponent mask after shift
    float32 o;

    o.u = (f.u & 0x7f) << 21;     // exponent/mantissa bits
    uint32_t exp = shifted_exp & o.u;   // just the exponent
    o.u += (127 - 2) << 23;        // exponent adjust

    // handle exponent special cases
    if (exp == shifted_exp) // Inf/NaN?
        o.u += (128 - 3) << 23;    // extra exp adjust
    else if (exp == 0) // Zero/Denormal?
    {
        o.u += 1 << 23;             // extra exp adjust
        o.f -= magic.f;             // renormalize
    }

    o.u |= (f.u & 0x80) << 24;    // sign bit
    return o;
}

#endif