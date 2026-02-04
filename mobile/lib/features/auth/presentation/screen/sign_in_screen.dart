import 'package:flutter/material.dart';
import 'package:mobile/core/theme/theme.dart';
import 'package:mobile/core/common/widget/card_general.dart';
import 'package:mobile/core/common/widget/primary_button.dart';

import '../../../../core/common/widget/input_text.dart';
import '../../../../core/utils/pref_helper.dart';
import '../../../../core/utils/size_extension.dart';

class SignInScreen extends StatefulWidget {
  const SignInScreen({super.key});

  @override
  State<SignInScreen> createState() => _SignInScreenState();
}

class _SignInScreenState extends State<SignInScreen> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;
  bool _rememberMe = false;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _rememberMe = PrefHelper.instance.isRememberMe;
    if (_rememberMe) {
      _emailController.text = PrefHelper.instance.email;
      _passwordController.text = PrefHelper.instance.password;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text("Welcome Back!", style: AppFont.medium16),
              Text("Please sign in to continue", style: AppFont.reguler14),
              CardGeneral(
                margin: EdgeInsets.symmetric(vertical: 24),
                child: Column(
                  children: [
                    InputText(
                      hintText: "Masukan alamat email anda",
                      title: "Email",
                      controller: _emailController,
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Email tidak boleh kosong';
                        }
                        if (!value.contains('@')) {
                          return 'Email tidak valid';
                        }
                        return null;
                      },
                    ),
                    widget.height(16),
                    InputText(
                      hintText: "Masukan password anda",
                      title: "Password",
                      controller: _passwordController,
                      obscureText: _obscurePassword,
                      // contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 0),
                      icon: InkWell(
                        child: Icon(
                          _obscurePassword ? Icons.visibility_off : Icons.visibility,
                          size: 16,
                        ),
                        onTap: () {
                          setState(() {
                            _obscurePassword = !_obscurePassword;
                          });
                        },
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Password tidak boleh kosong';
                        }
                        return null;
                      },
                    ),
                    widget.height(24),
                    PrimaryButton(title: "Login", onPressed: () {}),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
