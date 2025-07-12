import 'package:flutter/material.dart';

class PaymentScreen extends StatefulWidget {
  const PaymentScreen({super.key});

  @override
  State<PaymentScreen> createState() => _PaymentScreenState();
}

class _PaymentScreenState extends State<PaymentScreen> {
  final List<Map<String, dynamic>> _paymentMethods = [
    {
      'type': 'Credit Card',
      'icon': Icons.credit_card,
      'name': 'Visa ending in 1234',
      'isDefault': true,
    },
    {
      'type': 'Credit Card',
      'icon': Icons.credit_card,
      'name': 'Mastercard ending in 5678',
      'isDefault': false,
    },
    {
      'type': 'PayPal',
      'icon': Icons.account_balance_wallet,
      'name': 'john.doe@example.com',
      'isDefault': false,
    },
  ];

  int _selectedPaymentIndex = 0;

  @override
  void initState() {
    super.initState();
    // Find default payment method
    final defaultIndex = _paymentMethods.indexWhere((method) => method['isDefault'] == true);
    if (defaultIndex != -1) {
      _selectedPaymentIndex = defaultIndex;
    }
  }

  void _addNewPaymentMethod() {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => const _AddPaymentMethodScreen(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Payment Methods'),
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Your Payment Methods',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
              ),
              const SizedBox(height: 24),
              // Payment methods list
              Expanded(
                child: ListView.builder(
                  itemCount: _paymentMethods.length,
                  itemBuilder: (context, index) {
                    final method = _paymentMethods[index];
                    final isSelected = index == _selectedPaymentIndex;
                    
                    return Card(
                      margin: const EdgeInsets.only(bottom: 12),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                        side: BorderSide(
                          color: isSelected ? Colors.blue : Colors.transparent,
                          width: 2,
                        ),
                      ),
                      child: ListTile(
                        leading: Icon(
                          method['icon'],
                          color: isSelected ? Colors.blue : Colors.grey.shade700,
                        ),
                        title: Text(
                          method['name'],
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: isSelected ? Colors.blue : Colors.black,
                          ),
                        ),
                        subtitle: Text(method['type']),
                        trailing: method['isDefault']
                            ? Container(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 8,
                                  vertical: 4,
                                ),
                                decoration: BoxDecoration(
                                  color: Colors.blue.shade100,
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Text(
                                  'Default',
                                  style: TextStyle(
                                    color: Colors.blue.shade800,
                                    fontWeight: FontWeight.bold,
                                    fontSize: 12,
                                  ),
                                ),
                              )
                            : null,
                        onTap: () {
                          setState(() {
                            _selectedPaymentIndex = index;
                          });
                        },
                      ),
                    );
                  },
                ),
              ),
              // Add new payment method button
              SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: _addNewPaymentMethod,
                  icon: const Icon(Icons.add),
                  label: const Text('Add Payment Method'),
                ),
              ),
              const SizedBox(height: 16),
              // Set as default button
              if (!_paymentMethods[_selectedPaymentIndex]['isDefault'])
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () {
                      setState(() {
                        for (var i = 0; i < _paymentMethods.length; i++) {
                          _paymentMethods[i]['isDefault'] = i == _selectedPaymentIndex;
                        }
                      });
                      
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Default payment method updated'),
                          backgroundColor: Colors.green,
                        ),
                      );
                    },
                    child: const Text('Set as Default'),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class _AddPaymentMethodScreen extends StatefulWidget {
  const _AddPaymentMethodScreen();

  @override
  State<_AddPaymentMethodScreen> createState() => _AddPaymentMethodScreenState();
}

class _AddPaymentMethodScreenState extends State<_AddPaymentMethodScreen> {
  final _formKey = GlobalKey<FormState>();
  final _cardNumberController = TextEditingController();
  final _nameController = TextEditingController();
  final _expiryController = TextEditingController();
  final _cvvController = TextEditingController();
  bool _isLoading = false;

  @override
  void dispose() {
    _cardNumberController.dispose();
    _nameController.dispose();
    _expiryController.dispose();
    _cvvController.dispose();
    super.dispose();
  }

  void _savePaymentMethod() {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isLoading = true;
      });

      // Simulate API call
      Future.delayed(const Duration(seconds: 2), () {
        setState(() {
          _isLoading = false;
        });
        
        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Payment method added successfully'),
            backgroundColor: Colors.green,
          ),
        );
        
        // Navigate back
        Navigator.pop(context);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Add Payment Method'),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Card type icons
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.credit_card,
                      size: 40,
                      color: Colors.blue.shade700,
                    ),
                    const SizedBox(width: 16),
                    Icon(
                      Icons.credit_card,
                      size: 40,
                      color: Colors.red.shade700,
                    ),
                    const SizedBox(width: 16),
                    Icon(
                      Icons.credit_card,
                      size: 40,
                      color: Colors.amber.shade700,
                    ),
                  ],
                ),
                const SizedBox(height: 32),
                // Card number field
                TextFormField(
                  controller: _cardNumberController,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(
                    labelText: 'Card Number',
                    hintText: 'XXXX XXXX XXXX XXXX',
                    prefixIcon: Icon(Icons.credit_card),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter your card number';
                    }
                    if (value.replaceAll(' ', '').length != 16) {
                      return 'Card number must be 16 digits';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
                // Cardholder name field
                TextFormField(
                  controller: _nameController,
                  decoration: const InputDecoration(
                    labelText: 'Cardholder Name',
                    prefixIcon: Icon(Icons.person),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter the cardholder name';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
                // Expiry date and CVV fields
                Row(
                  children: [
                    // Expiry date field
                    Expanded(
                      child: TextFormField(
                        controller: _expiryController,
                        keyboardType: TextInputType.number,
                        decoration: const InputDecoration(
                          labelText: 'Expiry Date',
                          hintText: 'MM/YY',
                          prefixIcon: Icon(Icons.calendar_today),
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter expiry date';
                          }
                          if (!RegExp(r'^\d{2}/\d{2}$').hasMatch(value)) {
                            return 'Use format MM/YY';
                          }
                          return null;
                        },
                      ),
                    ),
                    const SizedBox(width: 16),
                    // CVV field
                    Expanded(
                      child: TextFormField(
                        controller: _cvvController,
                        keyboardType: TextInputType.number,
                        obscureText: true,
                        decoration: const InputDecoration(
                          labelText: 'CVV',
                          hintText: 'XXX',
                          prefixIcon: Icon(Icons.lock),
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter CVV';
                          }
                          if (value.length < 3 || value.length > 4) {
                            return 'CVV must be 3-4 digits';
                          }
                          return null;
                        },
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 32),
                // Save button
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _isLoading ? null : _savePaymentMethod,
                    child: _isLoading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.white,
                            ),
                          )
                        : const Text('Save Payment Method'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
