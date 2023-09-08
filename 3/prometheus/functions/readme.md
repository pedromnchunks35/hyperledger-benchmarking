# Functions
|Function|Description|
|--|--|
|abs(v instant-vector)|returns the input vector with all sample values converted to their absolute value|
|absent(v instant-vector)|returns an empty vector case there are elements in the vector, otherwise return vector with 1|
|absent_over_time(v range-vector)|returns an empty vector case there are elements in the vector, otherwise return vector with 1|
|ceil(v instant-vector)|rounds the sample values of all elements in **v** up to the nearest integer|
|changes(v range-vector)|returns the number of times its value has changed within the provided time range as an instant vector|
|clamp(v instant_vector,min scalar,max scalar)|clamps the samples values of all elements in **v** to have a lower limit of **min** and an upper limit of **max**,(case min>max,return empty vector. Case min or max is NaN return NaN)|
|clamp_max(v instant-vector,max scalar)|clamps the sample values of all elements in **v** to have an upper limit of **max**|
|clamp_min(v instant-vector,min scalar)|clamps the sample values of all elements in **v** to have a lower limit of **min**|
|day_of_month(v=vector(time()) instant-vector)|returns the day of the month for each of the given times in UTC. Returned values are from 1 to 31|
|day_of_year(v=vector(time()) instant vector)|returns the day of the year for each of the given times in UTC. Returned values are from 1 to 365 for non-leap years and 1 to 366 in leap years|
|days_in_month(v=vector(time()) instant-vector)|returns number of days in the month for each of the given times in UTC. Return values are from 28 to 31|
|delta(v range-vector)|calculates the difference between the first and last value of each time series element in a range vector **v**, returning an instant vector with the given deltas and equivalent labels. The delta is extrapolated to cover the full time range as specified in the range vector selector, so that it is possible to get a non-integer result even if the sample values are all integers|
|deriv(v range-vector)|Calculates the per-second derivative of the time series in a range vector v, using simple linear regression. The range vector must have at least two samples in order to perform the calculation. When +Inf or -Inf are found in the range vector, the slope and offset value calculated will be NaN. deriv should only be used with gauges.|
|exp(v instant-vector)|Calculates the exponential function for all elements in the instant-vector v. Special cases include:<br>- Exp(+Inf) = +Inf<br>- Exp(NaN) = NaN|
|floor(v instant-vector)|Rounds down the sample values of all elements in the instant-vector v to the nearest integer.|
|histogram_count(v instant-vector)|Returns the count of observations stored in a native histogram. Samples that are not native histograms are ignored and do not show up in the returned vector.|
|histogram_sum(v instant-vector)|Returns the sum of observations stored in a native histogram.|
|histogram_fraction(lower scalar, upper scalar, v instant-vector)|Returns the estimated fraction of observations between the provided lower and upper values for a native histogram. Samples that are not native histograms are ignored and do not show up in the returned vector.|
|histogram_quantile(φ scalar, b instant-vector)|Calculates the φ-quantile (0 ≤ φ ≤ 1) from either a conventional histogram or a native histogram. For a detailed explanation of φ-quantiles and the usage of the (conventional) histogram metric type, refer to histograms and summaries documentation.|
|histogram_stddev(v instant-vector)|Returns the estimated standard deviation of observations in a native histogram, calculated based on the geometric mean of the buckets where the observations lie. Samples that are not native histograms are ignored and do not appear in the returned vector.|
|histogram_stdvar(v instant-vector)|Returns the estimated standard variance of observations in a native histogram.|
|holt_winters(v range-vector, sf scalar, tf scalar)|Produces a smoothed value for a time series based on the range in v. The smoothing factor (sf) controls the importance given to old data, with lower values giving more importance to old data. The trend factor (tf) determines the consideration of trends in the data, with higher values considering trends more. Both sf and tf must be between 0 and 1.|
|hour(v=vector(time()) instant-vector)|Returns the hour of the day for each of the given times in UTC. Returned values range from 0 to 23.|
|idelta(v range-vector)|Calculates the difference between the last two samples in the range vector v and returns an instant vector with the calculated deltas and equivalent labels. It's important to note that idelta should only be used with gauges.|
|increase(v range-vector)|Calculates the increase in the time series represented by the range vector v. This calculation accounts for breaks in monotonicity, such as counter resets due to target restarts, and extrapolates the increase to cover the full time range as specified in the range vector selector. This means you can get a non-integer result even if a counter increases only by integer increments.|
|irate(v range-vector)|Calculates the per-second instant rate of increase for the time series represented by the range vector v. This calculation is based on the last two data points and automatically adjusts for breaks in monotonicity, such as counter resets due to target restarts.|
|label_join(v instant-vector, dst_label string, separator string, src_label_1 string, src_label_2 string, ...)|For each time series in the instant-vector v, this function joins all the values of the specified source labels (src_label_1, src_label_2, and so on) using the provided separator and returns the time series with the label dst_label containing the joined value. You can use any number of source labels in this function.|
|label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string)|This function matches the regular expression regex against the value of the label src_label in each time series of the instant-vector v. If a match is found, the value of the label dst_label in the returned time series will be the expansion of replacement, while the original labels from the input are retained. You can reference capturing groups in the regular expression using $1, $2, etc., or use $name for named capturing groups. If the regular expression doesn't match, the time series is returned unchanged.|
|ln(v instant-vector)|Calculates the natural logarithm for all elements in the instant-vector v. Special cases include: <br>- ln(+Inf) = +Inf<br>- ln(0) = -Inf<br>- ln(x < 0) = NaN<br>- ln(NaN) = NaN|
|log2(v instant-vector)|Calculates the binary logarithm for all elements in the instant-vector v. Special cases are as follows and are equivalent to those in ln|
|log10(v instant-vector)|Calculates the decimal logarithm for all elements in the instant-vector v. Special cases are equivalent to those in ln.|
|minute(v=vector(time()) instant-vector)|Returns the minute of the hour for each of the given times in UTC. Returned values range from 0 to 59.|
|month(v=vector(time()) instant-vector)|Returns the month of the year for each of the given times in UTC. Returned values range from 1 to 12, where 1 corresponds to January, and so on.|
|predict_linear(v range-vector, t scalar)|Predicts the value of the time series t seconds from now, based on the range vector v, using simple linear regression. To perform this calculation, the range vector must have at least two samples. If +Inf or -Inf values are found in the range vector, the calculated slope and offset values will be NaN.|
|rate(v range-vector)|Calculates the per-second average rate of increase for the time series represented by the range vector v. This calculation accounts for breaks in monotonicity, such as counter resets due to target restarts. Additionally, the calculation extrapolates to the ends of the time range, accommodating missed scrapes or imperfect alignment of scrape cycles with the range's time period.|
|resets(v range-vector)|For each input time series, the resets() function returns the number of counter resets within the provided time range as an instant vector. In this context:<br>- Any decrease in the value between two consecutive float samples is interpreted as a counter reset.<br>For native histograms, resets() detects resets in a more complex way:<br>- Any decrease in any bucket, including the zero bucket, or in the count of observations constitutes a counter reset.<br>- Additionally, resets() considers the disappearance of any previously populated bucket, an increase in bucket resolution, or a decrease in the zero-bucket width as counter resets.|
|round(v instant-vector, to_nearest=1 scalar)|Rounds the sample values of all elements in the instant-vector v to the nearest integer. Ties are resolved by rounding up. You can optionally specify the nearest multiple to which the sample values should be rounded using the to_nearest argument. This multiple can also be a fraction.|
|scalar(v instant-vector)|Given a single-element input vector, the scalar() function returns the sample value of that single element as a scalar. If the input vector does not have exactly one element, scalar() will return NaN.|
|sgn(v instant-vector)|Returns a vector with all sample values converted to their sign. The sign is defined as follows:<br>- 1 if v is positive,<br>- -1 if v is negative, and<br>- 0 if v is equal to zero.|
|sort(v instant-vector)|Returns vector elements sorted by their sample values in ascending order. For native histograms, sorting is based on the sum of observations.|
|sort_desc(v instant-vector)|Returns vector elements sorted by their sample values in descending order. For native histograms, sorting is based on the sum of observations.|
|sqrt(v instant-vector)|Calculates the square root of all elements in the instant-vector v.|
|time()|Returns the number of seconds since January 1, 1970 UTC. Please note that this function does not actually return the current time but rather the time at which the expression is to be evaluated.|
|timestamp(v instant-vector)|Returns the timestamp of each of the samples in the given vector as the number of seconds since January 1, 1970 UTC. This function also works with histogram samples.|
|vector(s scalar)|Returns the scalar s as a vector with no labels.|
|year(v=vector(time()) instant-vector)|Returns the year for each of the given times in UTC.|
|avg_over_time(range-vector)|Calculates the average value of all points in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results. All values in the interval have the same weight in the aggregation, regardless of spacing.
|min_over_time(range-vector)|Finds the minimum value of all points in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|max_over_time(range-vector)|Finds the maximum value of all points in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|sum_over_time(range-vector)|Calculates the sum of all values in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|count_over_time(range-vector)|Counts all values in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|quantile_over_time(scalar, range-vector)|Calculates the φ-quantile (0 ≤ φ ≤ 1) of the values in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|stddev_over_time(range-vector)|Calculates the population standard deviation of the values in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|stdvar_over_time(range-vector)|Calculates the population standard variance of the values in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|last_over_time(range-vector)|Returns the most recent point value in the specified time interval for each series in the range-vector and returns an instant vector with per-series aggregation results.
|present_over_time(range-vector)|Returns the value 1 for any series in the specified time interval for each series in the range-vector.|
|acos(v instant-vector)|Calculates the arccosine of all elements in v, working in radians (special cases).|
|acosh(v instant-vector)|Calculates the inverse hyperbolic cosine of all elements in v, working in radians (special cases).|
|asin(v instant-vector)|Calculates the arcsine of all elements in v, working in radians (special cases).|
|asinh(v instant-vector)|Calculates the inverse hyperbolic sine of all elements in v, working in radians (special cases).|
|atan(v instant-vector)|Calculates the arctangent of all elements in v, working in radians (special cases).|
|atanh(v instant-vector)|Calculates the inverse hyperbolic tangent of all elements in v, working in radians (special cases).|
|cos(v instant-vector)|Calculates the cosine of all elements in v, working in radians (special cases).|
|cosh(v instant-vector)|Calculates the hyperbolic cosine of all elements in v, working in radians (special cases).
|sin(v instant-vector)|Calculates the sine of all elements in v, working in radians (special cases).|
|sinh(v instant-vector)|Calculates the hyperbolic sine of all elements in v, working in radians (special cases).|
|tan(v instant-vector)|Calculates the tangent of all elements in v, working in radians (special cases).|
|tanh(v instant-vector)|Calculates the hyperbolic tangent of all elements in v, working in radians (special cases).|