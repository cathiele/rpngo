package convert

import (
	"fmt"
	"math"
)

//
// Implementation:
//
// Step 1: interpret the type string via a series of aliases
// Step 2: flip term 1 if needed, or return an error if there is no way to
//   convert
// Step 3: convert numerators
// Step 4: If needed
//  a) invert
//  b) convert numerators
//  c) invert
//
// Done!

type Unit struct {
	scale  float64
	offset float64
	names  []string
}

//
// Conversions to meters
//

var distantConvert = []Unit{
	{1.0000, 0, []string{"meter", "m", "meters", "metre", "metres"}},
	{1609.3440, 0, []string{"mile", "miles", "mi"}},
	{0.0254, 0, []string{"inch", "inches", "in"}},
	{0.0100, 0, []string{"centimeter", "centimeters", "cm"}},
	{0.1000, 0, []string{"decimeter", "decimeters"}},
	{0.0010, 0, []string{"millimeter", "millimeters", "mm"}},
	{1000.0000, 0, []string{"kilometer", "kilometers", "km"}},
	{0.3048, 0, []string{"foot", "feet", "ft"}},
	{0.9144, 0, []string{"yard", "yards", "yd"}},
	{0.155849127913, 0, []string{"_gallonmeters"}},
	{0.0779245639837, 0, []string{"_pintmeter"}},
	{0.0170183442005, 0, []string{"_teaspoonmeters"}},
	{0.0245446996271, 0, []string{"_tablespoonmeters"}},
	{0.098178798467, 0, []string{"_quartmeter"}},
	{63.6149072152, 0, []string{"_acremeter"}},
}

//
// Conversion to seconds
//

var timeConvert = []Unit{
	{1.0000, 0, []string{"second", "seconds", "sec", "s"}},
	{0.0010, 0, []string{"millisecond", "milliseconds", "ms"}},
	{1.0e-6, 0, []string{"microsecond", "microseconds", "us"}},
	{1.0e-9, 0, []string{"nanesecond", "nanoseconds", "ns"}},
	{60.0000, 0, []string{"minute", "minutes", "min"}},
	{3600.0000, 0, []string{"hour", "hours", "hr", "h"}},
	{86400.0000, 0, []string{"day", "days"}},
	{604800.0000, 0, []string{"week", "weeks"}},
}

//
// Weight/Mass Conversions (planet earth :) )
//

var massConvert = []Unit{
	{1.00000, 0, []string{"gram", "grams", "g"}},
	{1000.00000, 0, []string{"kilogram", "kilograms", "kg"}},
	{28.3495231, 0, []string{"ounce", "ounces", "oz"}},
	{453.59237, 0, []string{"pound", "pounds", "lbs", "lb"}},
	{907184.74000, 0, []string{"ton", "tons"}},
	{101.9680, 0, []string{"newton", "newtons", "n"}},
	{101968.0000, 0, []string{"kilonewton", "kilonewtons"}},
	{10196800.0000, 0, []string{"_barnewton"}},
}

//
// Cycles
//

var cycleConvert = []Unit{
	{1.0000, 0, []string{"cycle", "cycles"}},
	{1000000.0000, 0, []string{"kilacycle", "kilacycles"}},
	{1000000.0000, 0, []string{"megacycle", "megacycles"}},
	{1000000000.0000, 0, []string{"gigacycle", "gigacycles"}},
}

//
// Memory
//

var memoryConvert = []Unit{
	{1.0000, 0, []string{"bit", "bits"}},
	{1024.0000, 0, []string{"kilobit", "kilobits"}},
	{8.0000, 0, []string{"byte", "bytes"}},
	{8192.0000, 0, []string{"kilobyte", "kilobytes", "kb"}},
	{8388608.0000, 0, []string{"megabyte", "megabytes", "mb"}},
	{8589934592.0000, 0, []string{"gigabyte", "gigabytes", "gb"}},
	{8796093022208.0000, 0, []string{"terabyte", "terabytes", "tb"}},
}

//
// Angles
//

var angleConvert = []Unit{
	{1.0000, 0, []string{"radians", "rad"}},
	{math.Pi / 180.0, 0, []string{"degrees", "deg"}},
}

//
// Energy
//

var energyConvert = []Unit{
	{1.0000, 0, []string{"joules", "j"}},
	{1000.0000, 0, []string{"kilojoules", "kj"}},
	{1000000.0000, 0, []string{"megajoules", "mj"}},
	{3600000.0000, 0, []string{"kwh"}},
	{745.699872, 0, []string{"hps"}},
	{1055.0600, 0, []string{"btu"}},        // EC standard
	{105506000.0000, 0, []string{"therm"}}, // EC standard
	{4.2000, 0, []string{"calorie", "cal", "calories"}},
	{4200.0000, 0, []string{"kilocalorie", "kilocalories"}},
}

//
// Temperature
//

var temperatureConvert = []Unit{
	{1.0, 0, []string{"c", "celsius"}},
	{5.0 / 9.0, -32, []string{"f", "fahrenheit"}},
	{1.0, -273.15, []string{"k", "kelvin"}},
}

//
// Aliases to help things along
//

/*
ALIASES = {
  'acre': '_acremeter*_acremeter',
  'acres': 'acre',
  'bar': '_barnewton/meter*meter',
  'cadence': 'cycles/minute',
  'gallon': '_gallonmeters*_gallonmeters*_gallonmeters',
  'gallons': 'gallon',
  'ghz': 'gigacycles/second',
  'hp': 'hps/second',
  'hz': 'cycles/second',
  'khz': 'kilacycles/second',
  'kilowatt': 'kilojoules/second',
  'kilowatts': 'kilowatt',
  'kw': 'kilojoules/second',
  'liter': 'decimeter*decimeter*decimeter',
  'litre': 'liter',
  'liters': 'liter',
  'litres': 'liter',
  'milliliter': 'cm*cm*cm',
  'milliliters': 'milliliter',
  'ml': 'milliliter',
  'megawatt': 'megajoules/second',
  'megawatts': 'megawatt',
  'mw': 'megajoules/second',
  'mhz': 'megacycles/second',
  'mph': 'miles/hour',
  'pascal': 'newton/meter*meter',
  'pascals': 'pascal',
  'pa': 'pascal',
  'kilopascal': 'kilonewtons/meter*meter',
  'kilopascals': 'kilopascal',
  'kpa': 'kilopascal',
  'pint': '_pintmeter*_pintmeter*_pintmeter',
  'pints': 'pint',
  'psi': 'pounds/inch*inch',
  'quart': '_quartmeter*_quartmeter*_quartmeter',
  'quarts': 'quart',
  'rpm': 'cycles/minute',
  'tablespoon': '_tablespoonmeters*_tablespoonmeters*_tablespoonmeters',
  'tablespoons': 'tablespoon',
  'tsp': '_teaspoonmeters*_teaspoonmeters*_teaspoonmeters',
  'teaspoon': 'tsp',
  'teaspoons': 'tsp',
  'watt': 'joules/second',
  'watts': 'watt',
}

#
# Classes
#

class Error(Exception):
  def __init__(self, msg=None):
    Exception.__init__(self, msg)


class DuplicateKey(Error):
  def __init__(self, keyname):
    Error.__init__(self, keyname)


class IllegalConversionBetweenRatioAndScalar(Error):
  pass


class UnknownConversionType(Error):
  def __init__(self, type_name):
    Error.__init__(self, type_name)


class IncompatibleConversionTypes(Error):
  def __init__(self, source_type, target_type):
    Error.__init__(self, '%s -> %s' % (source_type, target_type))


class Conversion:

  def __init__(self):
    self.convert_dict = {}
    self._InsertKeys('Distance', DISTANCE_CONVERT)
    self._InsertKeys('Time', TIME_CONVERT)
    self._InsertKeys('Force/Weight/Mass (Planet Earth)', MASS_CONVERT)
    self._InsertKeys('Cycles', CYCLE_CONVERT)
    self._InsertKeys('Memory', MEMORY_CONVERT)
    self._InsertKeys('Angles', ANGLE_CONVERT)
    self._InsertKeys('Energy', ENERGY_CONVERT)
    self._InsertKeys('Temperature', TEMPERATURE_CONVERT)

  def Convert(self, value, value_type, target_type):

    # extract type and class information

    source = self._AnalyzeType(value_type)
    target = self._AnalyzeType(target_type)

    # check for ratio incompatibility

    if source.IsRatio() != target.IsRatio():
      raise IllegalConversionBetweenRatioAndScalar()

    # check for inversion eligibility

    if (target.IsRatio() and
        (source.numerator_name == target.denominator_name)):
      target.Invert()

    # check for numerator compatibility

    if source.numerator_name != target.numerator_name:
      raise IncompatibleConversionTypes(source.numerator_name,
                                        target.numerator_name)

    # check for denominator compatibility, if needed

    if (source.IsRatio() and
        (source.denominator_name != target.denominator_name)):
      raise IncompatibleConversionTypes(source.denominator_name,
                                        target.denominator_name)

    # scale the value by each numerator

    for snum in source.numerator:
      value = self._ScaleUp(value, snum.scale_factor)
    for tnum in target.numerator:
      value = self._ScaleDown(value, tnum.scale_factor)

    # if needed, scale by each denominator

    if source.IsRatio():
      for sden in source.denominator:
        value = self._ScaleDown(value, sden.scale_factor)
      for tden in target.denominator:
        value = self._ScaleUp(value, tden.scale_factor)
      if target.inverted:
        value = 1.0 / value

    return value

  def _ScaleUp(self, value, scale_factor):
    if isinstance(scale_factor, tuple):
      return (value + scale_factor[1]) * scale_factor[0]
    return value * scale_factor

  def _ScaleDown(self, value, scale_factor):
    if isinstance(scale_factor, tuple):
      return (value / scale_factor[0]) - scale_factor[1]
    return value / scale_factor

  def DumpHelp(self):

    classes = {}

    for conversion_name, conversion in self.convert_dict.items():
      if conversion.class_name not in classes:
        classes[conversion.class_name] = []
      classes[conversion.class_name].append(conversion_name)

    for class_name in sorted(classes):
      sys.stdout.write('\n%s:\n' % class_name)
      names = sorted(classes[class_name])
      names.reverse()
      sub_names = []
      while names:
        name = names.pop()
        if name.startswith('_'):
          continue
        sub_names.append(name)
        if len(sub_names) == 4:
          self._DumpColumns(sub_names)
          sub_names = []
      if sub_names:
        self._DumpColumns(sub_names)

    sys.stdout.write('\nUseful Aliases:\n')
    for alias_name in sorted(ALIASES):
      sys.stdout.write('  %-15s ->  %s\n' % (alias_name, ALIASES[alias_name]))

  def _DumpColumns(self, name_list):

    sys.stdout.write('  ')
    for name in name_list:
      sys.stdout.write('%-15s ' % name)
    sys.stdout.write('\n')

  def _AnalyzeType(self, type_str):

    numerator_type_list, denominator_type_list = self._AnalyzeTypeStr(type_str)

    self._CheckForAliases(numerator_type_list, denominator_type_list)
    self._CheckForAliases(denominator_type_list, numerator_type_list)

    numerator = []
    denominator = []

    for numerator_type in numerator_type_list:
      if numerator_type not in self.convert_dict:
        raise UnknownConversionType(numerator_type)
      numerator.append(self.convert_dict[numerator_type])

    for denominator_type in denominator_type_list:
      if denominator_type not in self.convert_dict:
        raise UnknownConversionType(denominator_type)
      denominator.append(self.convert_dict[denominator_type])

    return ConversionData(numerator, denominator)

  def _CheckForAliases(self, numerator, denominator):

    index = 0
    while index < len(numerator):
      if numerator[index] in ALIASES:
        alias = ALIASES[numerator[index]]
        del numerator[index]
        if '/' in alias:
          n, d = alias.split('/')
          numerator.extend(n.split('*'))
          denominator.extend(d.split('*'))
        else:
          numerator.extend(alias.split('*'))
        index = 0
      else:
        index += 1

    numerator.sort()
    denominator.sort()

  def _AnalyzeTypeStr(self, type_str):
    if '/' in type_str:
      num, den = type_str.split('/')
      numerator_type = sorted(num.split('*'))
      denominator_type = sorted(den.split('*'))
    else:
      numerator_type = sorted(type_str.split('*'))
      denominator_type = []
    return numerator_type, denominator_type


  def _InsertKeys(self, class_name, data):

    for conversion_tuple in data:
      scale_factor = conversion_tuple[0]
      keys = conversion_tuple[1:]
      for key in keys:
        if key in self.convert_dict:
          raise DuplicateKey(key)
        self.convert_dict[key] = ConversionType(class_name, scale_factor)


class ConversionData:

  def __init__(self, numerator, denominator):
    """Constructor.

    Args:
      numerator: ConversionType
      denominator: ConversionType for a ratio, None otherwise
    """

    self.numerator = numerator
    self.denominator = denominator
    self.inverted = False
    self.numerator_name = self._BuildClassName(
        self.numerator, self.denominator)
    self.denominator_name = self._BuildClassName(
        self.denominator, self.numerator)

  def IsRatio(self):
    return self.denominator is not None

  def Invert(self):
    self.numerator, self.denominator = self.denominator, self.numerator
    self.numerator_name, self.denominator_name = (
        self.denominator_name, self.numerator_name)
    self.inverted = not self.inverted

  def _BuildClassName(self, numerator, denominator):

    name_list = [x.class_name for x in numerator]
    for check_name in [x.class_name for x in denominator]:
      if check_name in name_list:
        name_list.remove(check_name)
    return '*'.join(sorted(name_list))

class ConversionType:

  def __init__(self, class_name, scale_factor):
    self.class_name = class_name
    self.scale_factor = scale_factor
*/

func Debugme() {
	fmt.Printf("%v\n", distantConvert)
}
