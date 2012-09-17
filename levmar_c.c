#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <float.h>

#include "/home/tony/src/levmar-2.6/levmar.h"

void func(double *p, double *x, int m, int n, void *data)
{
  callback_func(p,x,data);
}

void jacfunc(double *p, double *jac, int m, int n, void *data)
{
  callback_jacfunc(p,jac,data);
}

void levmar( double* ygiven, double* p, const int n, const int m, void* e ) {
  double opts[LM_OPTS_SZ], info[LM_INFO_SZ];

  // optimization control parameters; passing to levmar NULL instead of opts reverts to defaults
  opts[0]=LM_INIT_MU; opts[1]=1E-15; opts[2]=1E-15; opts[3]=1E-20;
  opts[4]=LM_DIFF_DELTA; // relevant only if the finite difference Jacobian version is used

  // invoke the optimization function
  dlevmar_der(func, jacfunc, p, ygiven, m, n, 1000, opts, info, NULL, NULL, e); // with analytic Jacobian
  //   dlevmar_dif(f1, p, x, m, n, 1000, opts, info, NULL, NULL, in); // without Jacobian
}
